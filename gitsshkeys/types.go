// Package gitsshkeys provides Git SSH key management for the Lattice Runtime API.
package gitsshkeys

import (
	"time"

	"github.com/google/uuid"
)

// GitSSHKey represents a user's Git SSH public key.
type GitSSHKey struct {
	UserID    uuid.UUID `json:"user_id" format:"uuid"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	PublicKey string    `json:"public_key"`
}
