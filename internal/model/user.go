package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role int

const (
	RoleApplicant Role = iota
	RoleEmployer
	RoleCurator
)

func (r Role) IsValid() bool {
	switch r {
	case RoleApplicant, RoleEmployer, RoleCurator:
		return true
	default:
		return false
	}
}

type User struct {
	ID uuid.UUID `gorm:"primaryKey;column:id"`

	Email        string `gorm:"uniqueIndex;type:varchar(100);notNull;column:email"`
	PasswordHash string `gorm:"type:varchar(72);notNull;column:password_hash"`
	Role         Role   `gorm:"type:smallint;notNull;column:role"`
	IsVerified   bool   `gorm:"default:false;column:is_verified"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
