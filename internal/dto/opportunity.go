package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type CreateOpportunityRequest struct {
	TagNames        []string              `json:"tag_names"`
	Title           string                `json:"title" validate:"required,min=5,max=60"`
	Description     string                `json:"description" validate:"required,min=10"`
	OpportunityType model.OpportunityType `json:"opportunity_type" validate:"required,oneof=job internship volunteer event course"`
	WorkFormat      model.WorkFormat      `json:"work_format" validate:"required,oneof=remote hybrid onsite"`
	LocationCity    string                `json:"location_city" validate:"required"`
	Latitude        float64               `json:"latitude"`
	Longitude       float64               `json:"longitude"`
	SalaryMin       int                   `json:"salary_min" validate:"gte=0"`
	SalaryMax       int                   `json:"salary_max" validate:"gte=0"`
	ExperienceLevel model.ExpirienseLevel `json:"experience_level" validate:"required,oneof=no_experience junior middle senior"`
	ExpiresAt       *time.Time            `json:"expires_at"`
	EventDateStart  *time.Time            `json:"event_date_start"`
	EventDateEnd    *time.Time            `json:"event_date_end"`
}

type UpdateOpportunityRequest struct {
	Title           *string                `json:"title" validate:"omitempty,min=5,max=60"`
	Description     *string                `json:"description" validate:"omitempty,min=10"`
	OpportunityType *model.OpportunityType `json:"opportunity_type" validate:"omitempty,oneof=job internship volunteer event course"`
	WorkFormat      *model.WorkFormat      `json:"work_format" validate:"omitempty,oneof=remote hybrid onsite"`
	LocationCity    *string                `json:"location_city"`
	Latitude        *float64               `json:"latitude"`
	Longitude       *float64               `json:"longitude"`
	SalaryMin       *int                   `json:"salary_min" validate:"omitempty,gte=0"`
	SalaryMax       *int                   `json:"salary_max" validate:"omitempty,gte=0"`
	ExperienceLevel *model.ExpirienseLevel `json:"experience_level" validate:"omitempty,oneof=no_experience junior middle senior"`
	ExpiresAt       *time.Time             `json:"expires_at"`
	EventDateStart  *time.Time             `json:"event_date_start"`
	EventDateEnd    *time.Time             `json:"event_date_end"`
}

type UpdateModerationStatusRequest struct {
	Status model.ModerationStatus `json:"status" validate:"required,oneof=pending approved rejected"`
}

type ListOpportunitiesRequest struct {
	EmployerID       *string                 `form:"employer_id"`
	CuratorID        *string                 `form:"curator_id"`
	OpportunityType  *model.OpportunityType  `form:"opportunity_type"`
	WorkFormat       *model.WorkFormat       `form:"work_format"`
	ExperienceLevel  *model.ExpirienseLevel  `form:"experience_level"`
	ModerationStatus *model.ModerationStatus `form:"moderation_status"`
	LocationCity     *string                 `form:"location_city"`
	TagIDs           []uuid.UUID             `form:"tag_ids"`
	MinSalary        *int                    `form:"min_salary"`
	MaxSalary        *int                    `form:"max_salary"`
	ExpiresAfter     *bool                   `form:"expires_after"`
	Limit            int                     `form:"limit" binding:"required,min=1,max=100"`
	Offset           int                     `form:"offset" binding:"min=0"`
}

type AddTagsRequest struct {
	TagIDs []uuid.UUID `json:"tag_ids" validate:"required,min=1"`
}

type RemoveTagsRequest struct {
	TagIDs []uuid.UUID `json:"tag_ids" validate:"required,min=1"`
}

type OpportunityResponse struct {
	ID               uuid.UUID              `json:"id"`
	EmployerID       *uuid.UUID             `json:"employer_id"`
	CuratorID        *uuid.UUID             `json:"curator_id"`
	Tags             []TagResponse          `json:"tags"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	OpportunityType  model.OpportunityType  `json:"opportunity_type"`
	WorkFormat       model.WorkFormat       `json:"work_format"`
	LocationCity     string                 `json:"location_city"`
	Latitude         float64                `json:"latitude"`
	Longitude        float64                `json:"longitude"`
	SalaryMin        int                    `json:"salary_min"`
	SalaryMax        int                    `json:"salary_max"`
	ExperienceLevel  model.ExpirienseLevel  `json:"experience_level"`
	ExpiresAt        *time.Time             `json:"expires_at"`
	EventDateStart   *time.Time             `json:"event_date_start"`
	EventDateEnd     *time.Time             `json:"event_date_end"`
	ModerationStatus model.ModerationStatus `json:"moderation_status"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type ListOpportunitiesResponse struct {
	Opportunities []OpportunityResponse `json:"opportunities"`
	Total         int64                 `json:"total"`
	Limit         int                   `json:"limit"`
	Offset        int                   `json:"offset"`
}
