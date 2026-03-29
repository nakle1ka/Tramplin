package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type Employer = model.Employer

type EmployerRepository interface {
	Create(ctx context.Context, employer *Employer) error
	List(ctx context.Context, filters EmployerListFilters) ([]*Employer, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Employer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Employer, error)
	Update(ctx context.Context, id uuid.UUID, employer map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type employerRepository struct {
	db *gorm.DB
}

func NewEmployerRepository(db *gorm.DB) EmployerRepository {
	return &employerRepository{db: db}
}

func (r *employerRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *employerRepository) Create(ctx context.Context, employer *Employer) error {
	err := r.getDB(ctx).Create(employer).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *employerRepository) GetByID(ctx context.Context, id uuid.UUID) (*Employer, error) {
	var employer Employer
	err := r.getDB(ctx).
		Preload("User").
		First(&employer, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &employer, nil
}

func (r *employerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Employer, error) {
	var employer Employer

	err := r.getDB(ctx).
		Preload("User").
		First(&employer, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &employer, nil
}

func (r *employerRepository) Update(ctx context.Context, id uuid.UUID, employer map[string]any) error {
	err := r.getDB(ctx).
		Model(&model.Employer{}).
		Where("id = ?", id).
		Updates(employer).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *employerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Delete(&Employer{}, "id = ?", id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

type EmployerListFilters struct {
	CompanyName    *string
	VerifiedStatus *model.VerificationStatus
	Limit          int
	Offset         int
}

func (r *employerRepository) List(ctx context.Context, filters EmployerListFilters) ([]*Employer, int64, error) {
	var employers []*Employer
	var total int64

	db := r.getDB(ctx).Model(&Employer{})

	if filters.CompanyName != nil {
		db = db.Where("company_name ILIKE ?", "%"+*filters.CompanyName+"%")
	}

	if filters.VerifiedStatus != nil {
		db = db.Where("verified_status = ?", *filters.VerifiedStatus)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	db = db.Limit(filters.Limit)

	if filters.Offset > 0 {
		db = db.Offset(filters.Offset)
	}

	err = db.
		Preload("User").
		Order("created_at DESC").
		Find(&employers).Error

	if err != nil {
		return nil, 0, err
	}

	return employers, total, nil
}
