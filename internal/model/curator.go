package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Curator struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;column:id"`
	UserID uuid.UUID `gorm:"type:uuid;unique;notNull;column:user_id"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	FullName     string `gorm:"type:varchar(150);notNull;column:full_name"`
	IsSuperAdmin bool   `gorm:"default:false;column:is_super_admin"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (c *Curator) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}
