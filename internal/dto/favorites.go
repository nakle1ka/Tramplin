package dto

import "github.com/google/uuid"

type OpportunityFavoritesResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	SalaryMin   int    `json:"salary_min"`
	SalaryMax   int    `json:"salary_max"`
}

type FavoriteResponse struct {
	ID            uuid.UUID                    `json:"id"`
	ApplicantID   uuid.UUID                    `json:"applicant_id"`
	OpportunityID uuid.UUID                    `json:"opportunity_id"`
	Opportunity   OpportunityFavoritesResponse `json:"opportunity"`
}

type FavoriteListResponse struct {
	Items []FavoriteResponse `json:"items"`
	Total int64              `json:"total"`
}

type CreateFavoriteRequest struct {
	OpportunityID uuid.UUID `json:"opportunity_id" binding:"required"`
}

type ListFavoritesRequest struct {
	ApplicantID   *string `form:"applicant_id"`
	OpportunityID *string `form:"opportunity_id"`
	Limit         int     `form:"limit" binding:"min=1,max=100"`
	Offset        int     `form:"offset,min:0"`
}
