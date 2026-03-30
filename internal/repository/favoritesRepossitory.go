package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type FavoritesRepository interface {
	Create(ctx context.Context, fovorites *model.Favorites) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Favorites, error)
	List(ctx context.Context, opts FavoritesListOptions) ([]model.Favorites, int64, error)
	Update(ctx context.Context, id uuid.UUID, fovorites map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type favoritesRepository struct {
	db *gorm.DB
}

func NewFavoritesRepository(db *gorm.DB) FavoritesRepository {
	return &favoritesRepository{db: db}
}

func (r *favoritesRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *favoritesRepository) Create(ctx context.Context, fovorites *model.Favorites) error {
	return r.getDB(ctx).Create(fovorites).Error
}

func (r *favoritesRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Favorites, error) {
	var fovorites model.Favorites
	err := r.getDB(ctx).
		Preload("Applicant").
		Preload("Opportunity").
		Where("id = ?", id).
		First(&fovorites).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &fovorites, nil
}

func (r *favoritesRepository) Update(ctx context.Context, id uuid.UUID, fovorites map[string]any) error {
	err := r.getDB(ctx).
		Model(&model.Favorites{}).
		Where("id = ?", id).
		Updates(fovorites).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *favoritesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Delete(&model.Favorites{}, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

type FavoritesListOptions struct {
	ApplicantID   *uuid.UUID
	OpportunityID *uuid.UUID
	Limit         int
	Offset        int
}

func (r *favoritesRepository) List(ctx context.Context, opts FavoritesListOptions) ([]model.Favorites, int64, error) {
	var fovoritess []model.Favorites

	query := r.getDB(ctx).Model(&model.Favorites{}).Preload("Applicant").Preload("Opportunity")
	if opts.ApplicantID != nil {
		query = query.Where("applicant_id = ?", *opts.ApplicantID)
	}
	if opts.OpportunityID != nil {
		query = query.Where("opportunity_id = ?", *opts.OpportunityID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("created_at DESC")
	query = query.Limit(opts.Limit)
	query = query.Offset(opts.Offset)

	err := query.Find(&fovoritess).Error
	return fovoritess, total, err
}
