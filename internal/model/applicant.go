package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Applicant struct {
	ID uuid.UUID `gorm:"primaryKey;column:id"`

	UserID uuid.UUID `gorm:"uniqueIndex;not null; column:user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	Tags []*Tag `gorm:"many2many:tag_applicants;constraint:OnDelete:CASCADE"`

	FirstName  string `gorm:"type:varchar(50);not null; column:first_name"`
	SecondName string `gorm:"type:varchar(50);not null; column:second_name"`
	LastName   string `gorm:"type:varchar(50);not null; column:last_name"`

	University     string `gorm:"type:varchar(150);column:university"`
	GraduationYear int    `gorm:"type:smallint;column:graduation_year"`

	About          string  `gorm:"type:text;column:about"`
	PrivacySetting Privacy `gorm:"type:smallint;default:1;column:privacy_setting"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (a *Applicant) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}

type Privacy int

const (
	PrivacyPrivate Privacy = iota
	PrivacyPublic
	PrivacyContacts
)

func (p Privacy) IsValid() bool {
	switch p {
	case PrivacyPrivate, PrivacyPublic, PrivacyContacts:
		return true
	default:
		return false
	}
}
