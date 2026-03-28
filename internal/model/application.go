package model

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID uuid.UUID `gorm:"primaryKey; column:id"`

	ApplicantID uuid.UUID `gorm:"index; uniqueIndex:applicants_opportunities; column:applicant_id"`
	Applicant   Applicant `gorm:"foreignKey:ApplicantID; references:ID; constraint:OnDelete:CASCADE"`

	OpportunityID uuid.UUID   `gorm:"index; uniqueIndex:applicants_opportunities; column:opportunity_id"`
	Opportunity   Opportunity `gorm:"foreignKey:OpportunityID; references:ID; constraint:OnDelete:CASCADE"`

	Status AcceptStatus `gorm:"column:status;default:0"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type AcceptStatus int

const (
	AcceptStatusPending AcceptStatus = iota
	AcceptStatusAccepted
	AcceptStatusRejected
)

func (s AcceptStatus) IsValid() bool {
	switch s {
	case AcceptStatusPending, AcceptStatusAccepted, AcceptStatusRejected:
		return true
	default:
		return false
	}
}
