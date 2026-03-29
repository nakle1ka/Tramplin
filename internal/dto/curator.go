package dto

import (
	"time"

	"github.com/google/uuid"
)

type CuratorResponse struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UpdateRequest struct {
	FullName *string `json:"full_name,omitempty"`
}
