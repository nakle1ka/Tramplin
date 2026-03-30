package dto

import "github.com/google/uuid"

type CreateRecommendationRequest struct {
	RecipientID   uuid.UUID `json:"recipient_id" binding:"required"`
	OpportunityID uuid.UUID `json:"opportunity_id" binding:"required"`
}

type ListRecommendationsQuery struct {
	Limit  int `form:"limit,default=10" binding:"min=1,max=100"`
	Offset int `form:"offset,default=0" binding:"min=0"`
}

type RecomendationResponse struct {
	ID uuid.UUID `json:"id"`

	SenderID uuid.UUID         `json:"sender_id"`
	Sender   RecomendationInfo `json:"sender"`

	Recipient   RecomendationInfo `json:"recipient"`
	RecipientID uuid.UUID         `json:"recipient_id"`

	OpportunityID uuid.UUID                    `json:"opportunity_id"`
	Opportunity   OpportunityRecomendationInfo `json:"opportunity"`
}

type RecomendationInfo struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
}

type OpportunityRecomendationInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
