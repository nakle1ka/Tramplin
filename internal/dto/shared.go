package dto

import "github.com/google/uuid"

type ErrorResponse struct {
	Error string `json:"error"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
