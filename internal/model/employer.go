package model

import (
	"time"

	"github.com/google/uuid"
)

type VerificationStatus int

const (
	StatusPending VerificationStatus = iota
	StatusVerified
	StatusRejected
)

type Employer struct {
	ID     uuid.UUID `gorm:"primaryKey;column:id"`
	UserID uuid.UUID `gorm:"uniqueIndex;not null;column:user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	CompanyName    string             `gorm:"type:varchar(150);not null;column:company_name"`
	INN            string             `gorm:"type:varchar(12);not null;column:inn"`
	Description    string             `gorm:"type:text;column:description"`
	Website        string             `gorm:"type:varchar(100);column:website"`
	VerifiedStatus VerificationStatus `gorm:"type:smallint;default:0;column:verified_status"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}
