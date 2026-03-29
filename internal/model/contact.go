package model

import (
	"time"

	"github.com/google/uuid"
)

type Contact struct {
	ID uuid.UUID `gorm:"primaryKey"`

	SenderID uuid.UUID `gorm:"index; not null"`
	Sender   Applicant `gorm:"foreignKey:SenderID; references:ID"`

	RecipientID uuid.UUID `gorm:"index; not null"`
	Recipient   Applicant `gorm:"foreignKey:RecipientID; references:ID"`

	Status ContactStatus `gorm:"index; not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type ContactStatus int

const (
	ContactStatusPending ContactStatus = iota
	ContactStatusAccepted
	ContactStatusRejected
)

func (cs ContactStatus) IsValid() bool {
	switch cs {
	case ContactStatusPending, ContactStatusAccepted, ContactStatusRejected:
		return true
	default:
		return false
	}
}
