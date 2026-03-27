package model

import (
	"time"

	"github.com/google/uuid"
)

type Curator struct {
	ID     uuid.UUID `gorm:"primaryKey;column:id"`
	UserID uuid.UUID `gorm:"uniqueIndex;not null;column:user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	FullName     string `gorm:"type:varchar(150);not null;column:full_name"`
	IsSuperAdmin bool   `gorm:"default:false;column:is_super_admin"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}
