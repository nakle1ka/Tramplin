package model

import (
	"time"

	"github.com/google/uuid"
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

	Email        string `gorm:"uniqueIndex;type:varchar(100);not null;column:email"`
	PasswordHash string `gorm:"type:varchar(72);not null;column:password_hash"`
	Role         Role   `gorm:"type:smallint;not null;column:role"`
	IsVerified   bool   `gorm:"default:false;column:is_verified"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}
