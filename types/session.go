package types

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an active agent session.
type Session struct {
	ID        uuid.UUID     `json:"id"`
	AgentID   uuid.UUID     `json:"agent_id"`
	UserID    uuid.UUID     `json:"user_id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Health    SessionHealth `json:"health"`
}

// SessionHealth contains session health status.
type SessionHealth struct {
	Healthy bool   `json:"healthy"`
	Reason  string `json:"reason,omitempty"`
}
