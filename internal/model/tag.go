package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID            uuid.UUID      `gorm:"primaryKey"`
	Name          string         `gorm:"unique; not null"`
	Applicants    []*Applicant   `gorm:"many2many:tag_applicants"`
	Opportunities []*Opportunity `gorm:"many2many:tag_opportunities"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
