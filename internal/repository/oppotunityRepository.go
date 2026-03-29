package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type Opportunity = model.Opportunity
type Tag = model.Tag

type OpportunityRepository interface {
	Create(ctx context.Context, opportunity *Opportunity) error
	Update(ctx context.Context, id uuid.UUID, updates map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error

	List(ctx context.Context, filter OpportunityFilter) ([]*Opportunity, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Opportunity, error)
	GetByEmployerID(ctx context.Context, employerID uuid.UUID) ([]*Opportunity, error)
	GetByCuratorID(ctx context.Context, curatorID uuid.UUID) ([]*Opportunity, error)

	AddTags(ctx context.Context, opportunityID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTags(ctx context.Context, opportunityID uuid.UUID, tagIDs []uuid.UUID) error
}

type OpportunityFilter struct {
	EmployerID       *uuid.UUID
	CuratorID        *uuid.UUID
	OpportunityType  *model.OpportunityType
	WorkFormat       *model.WorkFormat
	ExperienceLevel  *model.ExpirienseLevel
	ModerationStatus *model.ModerationStatus
	LocationCity     *string
	TagIDs           []uuid.UUID
	MinSalary        *int
	MaxSalary        *int
	ExpiresAfter     *bool
	Limit            int
	Offset           int
}

type opportunityRepository struct {
	db *gorm.DB
}

func NewOpportunityRepository(db *gorm.DB) OpportunityRepository {
	return &opportunityRepository{db: db}
}

func (r *opportunityRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *opportunityRepository) Create(ctx context.Context, opportunity *Opportunity) error {
	return r.getDB(ctx).Create(opportunity).Error
}

func (r *opportunityRepository) GetByID(ctx context.Context, id uuid.UUID) (*Opportunity, error) {
	var opportunity Opportunity
	err := r.getDB(ctx).
		Preload("Tags").
		Preload("Employer").
		Preload("Curator").
		First(&opportunity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &opportunity, nil
}

func (r *opportunityRepository) GetByEmployerID(ctx context.Context, employerID uuid.UUID) ([]*Opportunity, error) {
	var opportunities []*Opportunity
	err := r.getDB(ctx).
		Preload("Tags").
		Where("employer_id = ?", employerID).
		Order("created_at DESC").
		Find(&opportunities).Error
	if err != nil {
		return nil, err
	}
	return opportunities, nil
}

func (r *opportunityRepository) GetByCuratorID(ctx context.Context, curatorID uuid.UUID) ([]*Opportunity, error) {
	var opportunities []*Opportunity
	err := r.getDB(ctx).
		Preload("Tags").
		Where("curator_id = ?", curatorID).
		Order("created_at DESC").
		Find(&opportunities).Error
	if err != nil {
		return nil, err
	}
	return opportunities, nil
}

func (r *opportunityRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	result := r.getDB(ctx).
		Model(&Opportunity{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *opportunityRepository) AddTags(ctx context.Context, opportunityID uuid.UUID, tagIDs []uuid.UUID) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		var opportunity Opportunity
		if err := tx.First(&opportunity, "id = ?", opportunityID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if len(tagIDs) == 0 {
			return nil
		}

		var tags []Tag
		if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}

		return tx.Model(&opportunity).Association("Tags").Append(&tags)
	})
}

func (r *opportunityRepository) RemoveTags(ctx context.Context, opportunityID uuid.UUID, tagIDs []uuid.UUID) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		var opportunity Opportunity
		if err := tx.First(&opportunity, "id = ?", opportunityID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}

		if len(tagIDs) == 0 {
			return nil
		}

		var tags []Tag
		if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}

		return tx.Model(&opportunity).Association("Tags").Delete(&tags)
	})
}

func (r *opportunityRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).Delete(&Opportunity{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *opportunityRepository) List(ctx context.Context, filter OpportunityFilter) ([]*Opportunity, int64, error) {
	query := r.getDB(ctx).Model(&Opportunity{}).Preload("Tags")

	if filter.EmployerID != nil {
		query = query.Where("employer_id = ?", *filter.EmployerID)
	}
	if filter.CuratorID != nil {
		query = query.Where("curator_id = ?", *filter.CuratorID)
	}
	if filter.OpportunityType != nil {
		query = query.Where("opportunity_type = ?", *filter.OpportunityType)
	}
	if filter.WorkFormat != nil {
		query = query.Where("work_format = ?", *filter.WorkFormat)
	}
	if filter.ExperienceLevel != nil {
		query = query.Where("experience_level = ?", *filter.ExperienceLevel)
	}
	if filter.ModerationStatus != nil {
		query = query.Where("moderation_status = ?", *filter.ModerationStatus)
	}
	if filter.LocationCity != nil {
		query = query.Where("location_city = ?", *filter.LocationCity)
	}
	if filter.MinSalary != nil {
		query = query.Where("salary_min >= ? OR salary_max >= ?", *filter.MinSalary, *filter.MinSalary)
	}
	if filter.MaxSalary != nil {
		query = query.Where("salary_max <= ?", *filter.MaxSalary)
	}
	if filter.ExpiresAfter != nil {
		if *filter.ExpiresAfter {
			query = query.Where("expires_at IS NULL OR expires_at > NOW()")
		} else {
			query = query.Where("expires_at IS NOT NULL AND expires_at <= NOW()")
		}
	}

	if len(filter.TagIDs) > 0 {
		query = query.Joins("JOIN tag_opportunities ON tag_opportunities.opportunity_id = opportunities.id").
			Where("tag_opportunities.tag_id IN ?", filter.TagIDs).
			Group("opportunities.id").
			Having("COUNT(DISTINCT tag_opportunities.tag_id) = ?", len(filter.TagIDs))
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var opportunities []*Opportunity
	err := query.Order("created_at DESC").Find(&opportunities).Error
	if err != nil {
		return nil, 0, err
	}

	return opportunities, total, nil
}
