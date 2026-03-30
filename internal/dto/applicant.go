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
	PrivacySetting model.Privacy `json:"privacy_setting"`
	Tags           []TagResponse `json:"tags,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ApplicantInfo struct {
	ID         uuid.UUID `json:"id"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	LastName   string    `json:"last_name"`
}

type UpdateApplicantRequest struct {
	FirstName      *string        `json:"first_name,omitempty"`
	SecondName     *string        `json:"second_name,omitempty"`
	LastName       *string        `json:"last_name,omitempty"`
	University     *string        `json:"university,omitempty"`
	GraduationYear *int           `json:"graduation_year,omitempty"`
	WorkExpirience *string        `json:"work_experience,omitempty"`
	About          *string        `json:"about,omitempty"`
	PrivacySetting *model.Privacy `json:"privacy_setting,omitempty"`
}

type ListApplicantsRequest struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset int `form:"offset" binding:"omitempty,min=0"`
}

type ListApplicantsResponse struct {
	Applicants []ApplicantInfo `json:"applicants"`
	Total      int64           `json:"total"`
}

type TagsRequest struct {
	TagIDs []uuid.UUID `json:"tag_ids"`
}
