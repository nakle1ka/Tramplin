package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type Curator = model.Curator

type CuratorRepository interface {
	Create(ctx context.Context, curator *Curator) error
	GetByID(ctx context.Context, id uuid.UUID) (*Curator, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Curator, error)
	Update(ctx context.Context, id uuid.UUID, curator map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type curatorRepository struct {
	db *gorm.DB
}

func NewCuratorRepository(db *gorm.DB) CuratorRepository {
	return &curatorRepository{db: db}
}

func (r *curatorRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *curatorRepository) Create(ctx context.Context, curator *Curator) error {
	return r.getDB(ctx).Create(curator).Error
}

func (r *curatorRepository) GetByID(ctx context.Context, id uuid.UUID) (*Curator, error) {
	var curator Curator
	err := r.getDB(ctx).First(&curator, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &curator, nil
}

func (r *curatorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*Curator, error) {
	var curator Curator

	err := r.getDB(ctx).
		Preload("User").
		First(&curator, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &curator, nil
}

func (r *curatorRepository) Update(ctx context.Context, id uuid.UUID, curator map[string]any) error {
	err := r.getDB(ctx).
		Model(&model.Curator{}).
		Where("id = ?", id).
		Updates(curator).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}

func (r *curatorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.getDB(ctx).Delete(&Curator{}, "id = ?", id).Error
}
