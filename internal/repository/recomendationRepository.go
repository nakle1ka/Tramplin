package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type RecomendationRepository interface {
	Create(ctx context.Context, contact *model.Recomendation) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Recomendation, error)
	List(ctx context.Context, opts RecomendationListOptions) ([]model.Recomendation, error)
	Update(ctx context.Context, id uuid.UUID, contact map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type recomendationRepository struct {
	db *gorm.DB
}

func NewRecomendationRepository(db *gorm.DB) RecomendationRepository {
	return &recomendationRepository{db: db}
}

func (r *recomendationRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *recomendationRepository) Create(ctx context.Context, contact *model.Recomendation) error {
	return r.getDB(ctx).Create(contact).Error
}

func (r *recomendationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Recomendation, error) {
	var contact model.Recomendation
	err := r.getDB(ctx).
		Preload("Sender").
		Preload("Recipient").
		Where("id = ?", id).
		First(&contact).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &contact, nil
}

func (r *recomendationRepository) Update(ctx context.Context, id uuid.UUID, contact map[string]any) error {
	err := r.getDB(ctx).
		Model(&model.Recomendation{}).
		Where("id = ?", id).
		Updates(contact).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *recomendationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Delete(&model.Recomendation{}, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

type RecomendationListOptions struct {
	ApplicantID *uuid.UUID
	SenderID    *uuid.UUID
	RecipientID *uuid.UUID
	Limit       int
	Offset      int
}

func (r *recomendationRepository) List(ctx context.Context, opts RecomendationListOptions) ([]model.Recomendation, error) {
	var contacts []model.Recomendation

	query := r.getDB(ctx).Preload("Sender").Preload("Recipient").Preload("Opportunity")
	if opts.ApplicantID != nil {
		query = query.Where("sender_id = ? OR recipient_id = ?", *opts.ApplicantID, *opts.ApplicantID)
	}
	if opts.SenderID != nil {
		query = query.Where("sender_id = ?", *opts.SenderID)
	}
	if opts.RecipientID != nil {
		query = query.Where("recipient_id = ?", *opts.RecipientID)
	}

	query = query.Order("created_at DESC")
	query = query.Limit(opts.Limit)
	query = query.Offset(opts.Offset)

	err := query.Find(&contacts).Error
	return contacts, err
}
