package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Privacy int

const (
	PrivacyPrivate Privacy = iota
	PrivacyPublic
	PrivacyContacts
)

type Applicant struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID uuid.UUID `gorm:"type:uuid;unique;notNull;column:user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	FirstName  string `gorm:"type:varchar(50);notNull;column:first_name"`
	SecondName string `gorm:"type:varchar(50);notNull;column:second_name"`
	LastName   string `gorm:"type:varchar(50);notNull;column:last_name"`

	University     string `gorm:"type:varchar(150);column:university"`
	GraduationYear int    `gorm:"type:smallint;column:graduation_year"`

	About          string  `gorm:"type:text;column:about"`
	PrivacySetting Privacy `gorm:"type:smallint;default:0;column:privacy_setting"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (a *Applicant) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}
