package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type ApplicantResponse struct {
	ID             uuid.UUID     `json:"id"`
	Email          string        `json:"email"`
	UserID         uuid.UUID     `json:"user_id"`
	FirstName      string        `json:"first_name"`
	SecondName     string        `json:"second_name"`
	LastName       string        `json:"last_name"`
	University     string        `json:"university,omitempty"`
	GraduationYear int           `json:"graduation_year,omitempty"`
	About          string        `json:"about,omitempty"`
	PrivacySetting int           `json:"privacy_setting"`
	Tags           []TagResponse `json:"tags,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type TagResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type UpdateApplicantRequest struct {
	FirstName      *string        `json:"first_name,omitempty"`
	SecondName     *string        `json:"second_name,omitempty"`
	LastName       *string        `json:"last_name,omitempty"`
	University     *string        `json:"university,omitempty"`
	GraduationYear *int           `json:"graduation_year,omitempty"`
	About          *string        `json:"about,omitempty"`
	PrivacySetting *model.Privacy `json:"privacy_setting,omitempty"`
}

type TagsRequest struct {
	TagIDs []uuid.UUID `json:"tag_ids"`
}
