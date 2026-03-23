package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Opportunity struct {
	ID uuid.UUID `gorm:"primaryKey; column:id"`

	// one of: employer, curator
	EmployerID *uuid.UUID `gorm:"uniqueIndex; column:employer_id"`
	Employer   Employer   `gorm:"foreignKey:EmployerID; constraint:OnDelete:CASCADE"`

	CuratorID *uuid.UUID `gorm:"uniqueIndex; column:curator_id"`
	Curator   Curator    `gorm:"foreignKey:CuratorID; constraint:OnDelete:CASCADE"`

	Tags []*Tag `gorm:"many2many:tag_opportunities; constraint:OnDelete:CASCADE"`

	Title           string          `gorm:"type:varchar(60); not null; column:title"`
	Description     string          `gorm:"not null; column:description"`
	OpportunityType OpportunityType `gorm:"index; not null; column:opportunity_type"`
	WorkFormat      WorkFormat      `gorm:"index; not null; column:work_format"`

	LocationCity string  `gorm:"index; type:varchar(255); column:location_city"`
	Latitude     float64 `gorm:"type:decimal(10,8); column:latitude"`
	Longitude    float64 `gorm:"type:decimal(11,8); column:longitude"`

	SalaryMin       int             `gorm:"index; column:salary_min"`
	SalaryMax       int             `gorm:"column:salary_max"`
	ExperienceLevel ExpirienseLevel `gorm:"index; column:experience_level"`

	ModerationStatus string     `gorm:"column:moderation_status"`
	ExpiresAt        *time.Time `gorm:"column:expires_at"`
	EventDateStart   *time.Time `gorm:"column:event_date_start"`
	EventDateEnd     *time.Time `gorm:"column:event_date_end"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

func (o *Opportunity) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}

type ExpirienseLevel int

const (
	ExperienceLevelIntern ExpirienseLevel = iota
	ExperienceLevelJunior
	ExperienceLevelMiddle
	ExperienceLevelSenior
)

type ModerationStatus int

const (
	ModerationStatusPending ModerationStatus = iota
	ModerationStatusApproved
	ModerationStatusRejected
)

type OpportunityType int

const (
	OpportunityTypeInternship OpportunityType = iota
	OpportunityTypeVacancy
	OpportunityTypeCareerEvent
	OpportunityTypeMentoringProgram
)

type WorkFormat int

const (
	WorkFormatOffice WorkFormat = iota
	WorkFormatHybrid
	WorkFormatRemote
)
