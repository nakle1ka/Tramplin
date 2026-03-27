package model

import (
	"time"

	"github.com/google/uuid"
)

type Opportunity struct {
	ID uuid.UUID `gorm:"primaryKey; column:id"`

	// one of: employer, curator
	EmployerID *uuid.UUID `gorm:"index; column:employer_id"`
	Employer   Employer   `gorm:"foreignKey:EmployerID; references:ID; constraint:OnDelete:CASCADE"`

	CuratorID *uuid.UUID `gorm:"index; column:curator_id"`
	Curator   Curator    `gorm:"foreignKey:CuratorID; references:ID; constraint:OnDelete:CASCADE"`

	Tags []*Tag `gorm:"many2many:tag_opportunities; constraint:OnDelete:CASCADE"`

	Title           string          `gorm:"type:varchar(60); not null; column:title"`
	Description     string          `gorm:"not null; column:description"`
	OpportunityType OpportunityType `gorm:"index; not null; column:opportunity_type"`
	WorkFormat      WorkFormat      `gorm:"index; column:work_format"`

	LocationCity string  `gorm:"index; type:varchar(255); column:location_city"`
	Latitude     float64 `gorm:"type:decimal(10,8); column:latitude"`
	Longitude    float64 `gorm:"type:decimal(11,8); column:longitude"`

	SalaryMin       int             `gorm:"index; column:salary_min"`
	SalaryMax       int             `gorm:"column:salary_max"`
	ExperienceLevel ExpirienseLevel `gorm:"index; column:experience_level"`

	ModerationStatus ModerationStatus `gorm:"column:moderation_status; not null; default:0"`
	ExpiresAt        *time.Time       `gorm:"column:expires_at"`
	EventDateStart   *time.Time       `gorm:"column:event_date_start"`
	EventDateEnd     *time.Time       `gorm:"column:event_date_end"`

	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at"`
}

type ExpirienseLevel int

const (
	ExperienceLevelIntern ExpirienseLevel = iota
	ExperienceLevelJunior
	ExperienceLevelMiddle
	ExperienceLevelSenior
)

func (el ExpirienseLevel) IsValid() bool {
	switch el {
	case ExperienceLevelIntern,
		ExperienceLevelJunior,
		ExperienceLevelMiddle,
		ExperienceLevelSenior:
		return true
	default:
		return false
	}
}

type ModerationStatus int

const (
	ModerationStatusPending ModerationStatus = iota
	ModerationStatusApproved
	ModerationStatusRejected
)

func (ms ModerationStatus) IsValid() bool {
	switch ms {
	case ModerationStatusApproved,
		ModerationStatusPending,
		ModerationStatusRejected:
		return true
	default:
		return false
	}
}

type OpportunityType int

const (
	OpportunityTypeInternship OpportunityType = iota
	OpportunityTypeVacancy
	OpportunityTypeCareerEvent
	OpportunityTypeMentoringProgram
)

func (ot OpportunityType) IsValid() bool {
	switch ot {
	case OpportunityTypeInternship,
		OpportunityTypeVacancy,
		OpportunityTypeCareerEvent,
		OpportunityTypeMentoringProgram:
		return true
	default:
		return false
	}
}

type WorkFormat int

const (
	WorkFormatOffice WorkFormat = iota
	WorkFormatHybrid
	WorkFormatRemote
)

func (wf WorkFormat) IsValid() bool {
	switch wf {
	case WorkFormatOffice,
		WorkFormatHybrid,
		WorkFormatRemote:
		return true
	default:
		return false
	}
}
