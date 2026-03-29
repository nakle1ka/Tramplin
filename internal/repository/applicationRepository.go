package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(ctx context.Context, application *model.Application) error
	UpdateAcceptStatus(ctx context.Context, id uuid.UUID, status model.AcceptStatus) error
	Delete(ctx context.Context, id uuid.UUID) error

	GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error)
	GetByOpportunityID(ctx context.Context, opportunityID uuid.UUID, limit, offset int) ([]*model.Application, int64, error)
	GetByApplicantID(ctx context.Context, applicantID uuid.UUID, limit, offset int) ([]*model.Application, int64, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *applicationRepository) Create(ctx context.Context, application *model.Application) error {
	return r.getDB(ctx).Create(application).Error
}

func (r *applicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	var application model.Application
	err := r.getDB(ctx).
		Preload("Applicant").
		Preload("Opportunity").
		Where("id = ?", id).
		First(&application).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &application, nil
}

func (r *applicationRepository) GetByOpportunityID(ctx context.Context, opportunityID uuid.UUID, limit, offset int) ([]*model.Application, int64, error) {
	var applications []*model.Application
	var total int64

	db := r.getDB(ctx)
	if err := db.Model(&model.Application{}).
		Where("opportunity_id = ?", opportunityID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Applicant").
		Preload("Opportunity").
		Where("opportunity_id = ?", opportunityID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&applications).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*model.Application{}, 0, ErrNotFound
		}
		return []*model.Application{}, 0, err
	}

	return applications, total, err
}

func (r *applicationRepository) GetByApplicantID(ctx context.Context, applicantID uuid.UUID, limit, offset int) ([]*model.Application, int64, error) {
	var applications []*model.Application
	var total int64

	db := r.getDB(ctx)

	if err := db.Model(&model.Application{}).
		Where("applicant_id = ?", applicantID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.
		Preload("Applicant").
		Preload("Opportunity").
		Where("applicant_id = ?", applicantID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&applications).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []*model.Application{}, 0, ErrNotFound
		}
		return []*model.Application{}, 0, err
	}

	return applications, total, err
}

func (r *applicationRepository) UpdateAcceptStatus(ctx context.Context, id uuid.UUID, status model.AcceptStatus) error {
	err := r.getDB(ctx).
		Model(&model.Application{}).
		Where("id = ?", id).
		Update("status", status).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *applicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Delete(&model.Application{}, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}
