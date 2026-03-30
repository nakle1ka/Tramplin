package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type ApplicationResponse struct {
	ID            uuid.UUID                   `json:"id"`
	OpportunityID uuid.UUID                   `json:"opportunity_id"`
	Opportunity   *OpportunityApplicationInfo `json:"opportunity,omitempty"`
	ApplicantID   uuid.UUID                   `json:"applicant_id"`
	Applicant     *ApplicationApplicantInfo   `json:"applicant,omitempty"`
	Status        model.AcceptStatus          `json:"status"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
}

type OpportunityApplicationInfo struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	LocationCity string    `json:"location_city"`
}

type ApplicationApplicantInfo struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}

type GetApplicationsRequest struct {
	OpportunityID *string `form:"opportunity_id"`
	ApplicantID   *string `form:"applicant_id"`
	Limit         int     `form:"limit" binding:"min=1,max=100"`
	Offset        int     `form:"offset" binding:"min=0"`
}

type CreateApplicationRequest struct {
	OpportunityID uuid.UUID `json:"opportunity_id" binding:"required"`
}

type UpdateApplicationStatusRequest struct {
	Status model.AcceptStatus `json:"status" binding:"required"`
}

type GetApplicationsByOpportunityRequest struct {
	OpportunityID uuid.UUID `form:"opportunity_id" binding:"required"`
	Limit         int       `form:"limit,default=20" binding:"min=1,max=100"`
	Offset        int       `form:"offset,default=0" binding:"min=0"`
}

type GetApplicationsByApplicantRequest struct {
	ApplicantID uuid.UUID `form:"applicant_id" binding:"required"`
	Limit       int       `form:"limit,default=20" binding:"min=1,max=100"`
	Offset      int       `form:"offset,default=0" binding:"min=0"`
}

type ApplicationsListResponse struct {
	Applications []*ApplicationResponse `json:"applications"`
	Total        int64                  `json:"total"`
	Limit        int                    `json:"limit"`
	Offset       int                    `json:"offset"`
}
