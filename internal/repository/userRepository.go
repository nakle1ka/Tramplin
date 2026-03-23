package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type User = model.User

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	SetVerify(ctx context.Context, id uuid.UUID, value bool) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	return r.getDB(ctx).Create(user).Error
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.getDB(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := r.getDB(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SetVerify(ctx context.Context, id uuid.UUID, value bool) error {
	return r.getDB(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Update("is_verified", value).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.getDB(ctx).Delete(&User{}, "id = ?", id).Error
}
