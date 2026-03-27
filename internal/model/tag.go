package model

import (
	"github.com/google/uuid"
)

type Tag struct {
	ID            uuid.UUID      `gorm:"primaryKey"`
	Name          string         `gorm:"unique; not null"`
	Applicants    []*Applicant   `gorm:"many2many:tag_applicants;constraint:OnDelete:CASCADE"`
	Opportunities []*Opportunity `gorm:"many2many:tag_opportunities;constraint:OnDelete:CASCADE"`
}
