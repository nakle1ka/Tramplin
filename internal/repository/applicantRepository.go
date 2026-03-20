package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type Applicant = model.Applicant

type ApplicantRepository interface {
	Create(ctx context.Context, applicant *Applicant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Applicant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Applicant, error)
	Update(ctx context.Context, applicant *Applicant) error
	Delete(ctx context.Context, id uuid.UUID) error
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

func (r *applicantRepository) Create(ctx context.Context, applicant *Applicant) error {
	return r.getDB(ctx).Create(applicant).Error
}

func (r *applicantRepository) GetByID(ctx context.Context, id uuid.UUID) (*Applicant, error) {
	var applicant Applicant
	err := r.getDB(ctx).First(&applicant, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &applicant, nil
}

func (r *applicantRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Applicant, error) {
	var applicant Applicant
	err := r.getDB(ctx).First(&applicant, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &applicant, nil
}

func (r *applicantRepository) Update(ctx context.Context, applicant *Applicant) error {
	return r.getDB(ctx).Save(applicant).Error
}

func (r *applicantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.getDB(ctx).Delete(&Applicant{}, "id = ?", id).Error
}
