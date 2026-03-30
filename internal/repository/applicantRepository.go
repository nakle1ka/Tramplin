package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type ApplicantRepository interface {
	Create(ctx context.Context, applicant *model.Applicant) error
	Update(ctx context.Context, id uuid.UUID, applicant map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error

	List(ctx context.Context, limit, offset int) ([]*model.Applicant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Applicant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Applicant, error)

	AddTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error
	GetTags(ctx context.Context, applicantID uuid.UUID) ([]*model.Tag, error)
	SetTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error
}

type applicantRepository struct {
	db *gorm.DB
}

func NewApplicantRepository(db *gorm.DB) ApplicantRepository {
	return &applicantRepository{db: db}
}

func (r *applicantRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *applicantRepository) Create(ctx context.Context, applicant *model.Applicant) error {
	return r.getDB(ctx).Create(applicant).Error
}

func (r *applicantRepository) List(ctx context.Context, limit, offset int) ([]*model.Applicant, error) {
	var applicants []*model.Applicant

	err := r.getDB(ctx).
		Preload("Tags").
		Preload("User").
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&applicants).Error

	if err != nil {
		return nil, err
	}

	return applicants, nil
}

func (r *applicantRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Applicant, error) {
	var applicant model.Applicant
	err := r.getDB(ctx).
		Preload("Tags").
		Preload("User").
		First(&applicant, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &applicant, nil
}

func (r *applicantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*model.Applicant, error) {
	var applicant model.Applicant
	err := r.getDB(ctx).
		Preload("Tags").
		Preload("User").
		First(&applicant, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &applicant, nil
}

func (r *applicantRepository) Update(ctx context.Context, id uuid.UUID, applicant map[string]any) error {
	result := r.getDB(ctx).Model(&model.Applicant{}).
		Where("id = ?", id).
		Updates(applicant)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *applicantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.getDB(ctx).Delete(&model.Applicant{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *applicantRepository) AddTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error {
	var applicant model.Applicant
	if err := r.getDB(ctx).First(&applicant, "id = ?", applicantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	var tags []*model.Tag
	if err := r.getDB(ctx).Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}

	return r.getDB(ctx).Model(&applicant).Association("Tags").Append(tags)
}

func (r *applicantRepository) RemoveTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error {
	var applicant model.Applicant
	if err := r.getDB(ctx).First(&applicant, "id = ?", applicantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	var tags []*model.Tag
	if err := r.getDB(ctx).Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
		return err
	}

	return r.getDB(ctx).Model(&applicant).Association("Tags").Delete(tags)
}

func (r *applicantRepository) GetTags(ctx context.Context, applicantID uuid.UUID) ([]*model.Tag, error) {
	var applicant model.Applicant
	if err := r.getDB(ctx).First(&applicant, "id = ?", applicantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var tags []*model.Tag
	err := r.getDB(ctx).Model(&applicant).Association("Tags").Find(&tags)
	return tags, err
}

func (r *applicantRepository) SetTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error {
	var applicant model.Applicant
	if err := r.getDB(ctx).First(&applicant, "id = ?", applicantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	var tags []*model.Tag
	if len(tagIDs) > 0 {
		if err := r.getDB(ctx).Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			return err
		}
	}

	return r.getDB(ctx).Model(&applicant).Association("Tags").Replace(tags)
}
