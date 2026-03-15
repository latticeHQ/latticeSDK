// Package licenses provides license management for the Lattice Runtime API.
package licenses

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// LicenseExpiryClaim is the JWT claim key for license expiration.
	LicenseExpiryClaim = "license_expires"
)

// License represents a Lattice license with JWT claims.
type License struct {
	ID         int32                  `json:"id"`
	UUID       uuid.UUID              `json:"uuid" format:"uuid"`
	UploadedAt time.Time              `json:"uploaded_at" format:"date-time"`
	// Claims are the JWT claims asserted by the license. A generic string map
	// is used to ensure all data from the server is parsed verbatim.
	Claims map[string]interface{} `json:"claims"`
}

// ExpiresAt returns the expiration time of the license.
// If the claim is missing or has an unexpected type, an error is returned.
func (l *License) ExpiresAt() (time.Time, error) {
	expClaim, ok := l.Claims[LicenseExpiryClaim]
	if !ok {
		return time.Time{}, fmt.Errorf("license_expires claim is missing")
	}

	if unix, ok := expClaim.(json.Number); ok {
		i64, err := unix.Int64()
		if err != nil {
			return time.Time{}, fmt.Errorf("license_expires claim is not a valid unix timestamp: %w", err)
		}
		return time.Unix(i64, 0), nil
	}

	return time.Time{}, fmt.Errorf("license_expires claim has unexpected type %T", expClaim)
}

// Trial returns whether the license is a trial license.
func (l *License) Trial() bool {
	if trial, ok := l.Claims["trail"].(bool); ok {
		return trial
	}
	return false
}

// AllFeaturesClaim returns whether the license grants all features.
func (l *License) AllFeaturesClaim() bool {
	if all, ok := l.Claims["all_features"].(bool); ok {
		return all
	}
	return false
}

// FeaturesClaims returns the explicit feature claims in the license.
// If checking for actual usage, also check AllFeaturesClaim.
func (l *License) FeaturesClaims() (map[string]int64, error) {
	strMap, ok := l.Claims["features"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("features key is unexpected type")
	}
	fMap := make(map[string]int64)
	for k, v := range strMap {
		jn, ok := v.(json.Number)
		if !ok {
			return nil, fmt.Errorf("feature %q has unexpected type", k)
		}
		n, err := jn.Int64()
		if err != nil {
			return nil, err
		}
		fMap[k] = n
	}
	return fMap, nil
}

// AddLicenseRequest is the request body for adding a new license.
type AddLicenseRequest struct {
	License string `json:"license" validate:"required"`
}
