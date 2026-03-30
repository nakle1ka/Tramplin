package model

import (
	"time"

	"github.com/google/uuid"
)

type Favorites struct {
	ID uuid.UUID `gorm:"primaryKey"`

	ApplicantID uuid.UUID `gorm:"index; uniqueIndex:favorites_idx; column:applicant_id"`
	Applicant   Applicant `gorm:"foreignKey:ApplicantID; references:ID"`

	OpportunityID uuid.UUID   `gorm:"index; uniqueIndex:favorites_idx; column:opportunity_id"`
	Opportunity   Opportunity `gorm:"foreignKey:OpportunityID; references:ID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
