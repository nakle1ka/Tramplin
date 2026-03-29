package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"gorm.io/gorm"
)

type ContactRepository interface {
	Create(ctx context.Context, contact *model.Contact) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Contact, error)
	List(ctx context.Context, opts ListOptions) ([]model.Contact, error)
	Update(ctx context.Context, id uuid.UUID, contact map[string]any) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type contactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) getDB(ctx context.Context) *gorm.DB {
	if txWrapper, ok := ctx.Value(ctxKey{}).(*Transaction); ok && txWrapper.Tx != nil {
		return txWrapper.Tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *contactRepository) Create(ctx context.Context, contact *model.Contact) error {
	return r.getDB(ctx).Create(contact).Error
}

func (r *contactRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Contact, error) {
	var contact model.Contact
	err := r.getDB(ctx).
		Preload("Sender").
		Preload("Recipient").
		Where("id = ?", id).
		First(&contact).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContactNotFound
		}
		return nil, err
	}
	return &contact, nil
}

func (r *contactRepository) Update(ctx context.Context, id uuid.UUID, contact map[string]any) error {
	err := r.getDB(ctx).
		Model(&model.Contact{}).
		Where("id = ?", id).
		Updates(contact).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContactNotFound
		}
		return err
	}

	return nil
}

func (r *contactRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.getDB(ctx).Delete(&model.Contact{}, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContactNotFound
		}
		return err
	}
	return nil
}

type ListOptions struct {
	ApplicantID *uuid.UUID
	Status      *model.ContactStatus
	ContactID   *uuid.UUID
	SenderID    *uuid.UUID
	RecipientID *uuid.UUID
	Limit       int
	Offset      int
}

func (r *contactRepository) List(ctx context.Context, opts ListOptions) ([]model.Contact, error) {
	var contacts []model.Contact

	query := r.getDB(ctx).Preload("Sender").Preload("Recipient")
	if opts.ApplicantID != nil {
		query = query.Where("sender_id = ? OR recipient_id = ?", *opts.ApplicantID, *opts.ApplicantID)
	}
	if opts.Status != nil {
		query = query.Where("status = ?", *opts.Status)
	}
	if opts.ContactID != nil {
		query = query.Where("id = ?", *opts.ContactID)
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
