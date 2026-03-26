package dto

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type RegisterRequest struct {
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	Role     model.Role `json:"role"`

	Applicant *ApplicantData `json:"applicant,omitempty"`
	Employer  *EmployerData  `json:"employer,omitempty"`
}

type ApplicantData struct {
	FirstName  string `json:"first_name" binding:"required"`
	SecondName string `json:"second_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
}

type EmployerData struct {
	CompanyName string `json:"company_name" binding:"required"`
	INN         string `json:"inn" binding:"required"`
	Description string `json:"description,omitempty"`
	Website     string `json:"website,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken string     `json:"access_token"`
	UserID      uuid.UUID  `json:"user_id"`
	Role        model.Role `json:"role"`
	IsVerified  bool       `json:"is_verified"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
