package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type UpdateEmployerFields struct {
	CompanyName *string `json:"company_name"`
	Description *string `json:"description"`
	Website     *string `json:"website"`
}

type UpdateVerificationStatusRequest struct {
	Status *model.VerificationStatus `json:"status" binding:"required"`
}

type EmployerResponse struct {
	ID             uuid.UUID                `json:"id"`
	UserID         uuid.UUID                `json:"user_id"`
	Email          string                   `json:"email"`
	CompanyName    string                   `json:"company_name"`
	INN            string                   `json:"inn"`
	Description    string                   `json:"description"`
	Website        string                   `json:"website"`
	VerifiedStatus model.VerificationStatus `json:"verified_status"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
}
