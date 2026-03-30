package model

import (
	"time"

	"github.com/google/uuid"
)

type Recomendation struct {
	ID uuid.UUID `gorm:"primaryKey"`

	SenderID uuid.UUID `gorm:"index; uniqueIndex:recomendation_index; not null; column:sender_id"`
	Sender   Applicant `gorm:"foreignKey:SenderID; references:ID"`

	RecipientID uuid.UUID `gorm:"index; uniqueIndex:recomendation_index; not null; column:recipient_id"`
	Recipient   Applicant `gorm:"foreignKey:RecipientID; references:ID"`

	OpportunityID uuid.UUID   `gorm:"index; uniqueIndex:recomendation_index; not null; column:opportunity_id"`
	Opportunity   Opportunity `gorm:"foreignKey:OpportunityID; references:ID"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
