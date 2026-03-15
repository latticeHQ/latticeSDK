package client

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

// Pagination sets pagination options for list endpoints.
type Pagination struct {
	// AfterID returns results after the given UUID.
	// Use this as an alternative to Offset for cursor-based pagination.
	AfterID uuid.UUID `json:"after_id,omitempty"`
	// Limit sets the maximum number of results per page.
	// If <= 0, no limit is applied.
	Limit int `json:"limit,omitempty"`
	// Offset is the number of results to skip.
	// Use offset=limit*page for page-based pagination.
	Offset int `json:"offset,omitempty"`
}

// AsRequestOption returns a RequestOption that applies pagination query parameters.
func (p Pagination) AsRequestOption() RequestOption {
	return func(r *http.Request) {
		q := r.URL.Query()
		if p.AfterID != uuid.Nil {
			q.Set("after_id", p.AfterID.String())
		}
		if p.Limit > 0 {
			q.Set("limit", strconv.Itoa(p.Limit))
		}
		if p.Offset > 0 {
			q.Set("offset", strconv.Itoa(p.Offset))
		}
		r.URL.RawQuery = q.Encode()
	}
}
