package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
	"gorm.io/gorm"
)

type ContactService interface {
	Create(ctx context.Context, dto CreateContactDTO) error

	ListFriends(ctx context.Context, dto ListFriendsDTO) ([]model.Contact, error)
	ListSentRequests(ctx context.Context, dto ListSentRequestsDTO) ([]model.Contact, error)
	ListReceivedRequests(ctx context.Context, dto ListReceivedRequestsDTO) ([]model.Contact, error)

	UpdateStatus(ctx context.Context, dto UpdateStatusDTO) error
	Delete(ctx context.Context, dto DeleteContactDTO) error
}

type contactService struct {
	repo          repository.ContactRepository
	applicantRepo repository.ApplicantRepository
}

func NewContactService(repo repository.ContactRepository, applicantRepo repository.ApplicantRepository) ContactService {
	return &contactService{
		repo:          repo,
		applicantRepo: applicantRepo,
	}
}

type CreateContactDTO struct {
	Auth        AuthContext
	RecipientID uuid.UUID
}

func (s *contactService) Create(ctx context.Context, dto CreateContactDTO) error {
	if dto.Auth.Role != model.RoleApplicant {
		return ErrForbidden
	}

	sender, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return err
	}
	if sender.ID == dto.RecipientID {
		return ErrInvalidInput
	}

	recipient, err := s.applicantRepo.GetByID(ctx, dto.RecipientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidInput
		}
		return err
	}

	contact := &model.Contact{
		ID:          uuid.New(),
		SenderID:    sender.ID,
		RecipientID: recipient.ID,
		Status:      model.ContactStatusPending,
	}

	return s.repo.Create(ctx, contact)
}

type ListFriendsDTO struct {
	Auth   AuthContext
	Limit  int
	Offset int
}

func (s *contactService) ListFriends(ctx context.Context, dto ListFriendsDTO) ([]model.Contact, error) {
	if dto.Auth.Role != model.RoleApplicant {
		return nil, ErrForbidden
	}
	if dto.Limit < 0 || dto.Limit > 100 || dto.Offset < 0 {
		return nil, ErrInvalidInput
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return nil, err
	}

	status := model.ContactStatusAccepted
	opts := repository.ContantListOptions{
		ApplicantID: &applicant.ID,
		Status:      &status,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
	}

	return s.repo.List(ctx, opts)
}

type ListSentRequestsDTO struct {
	Auth   AuthContext
	Limit  int
	Offset int
}

func (s *contactService) ListSentRequests(ctx context.Context, dto ListSentRequestsDTO) ([]model.Contact, error) {
	if dto.Auth.Role != model.RoleApplicant {
		return nil, ErrForbidden
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return nil, err
	}

	status := model.ContactStatusPending
	opts := repository.ContantListOptions{
		SenderID: &applicant.ID,
		Status:   &status,
		Limit:    dto.Limit,
		Offset:   dto.Offset,
	}

	return s.repo.List(ctx, opts)
}

type ListReceivedRequestsDTO struct {
	Auth   AuthContext
	Limit  int
	Offset int
}

func (s *contactService) ListReceivedRequests(ctx context.Context, dto ListReceivedRequestsDTO) ([]model.Contact, error) {
	if dto.Auth.Role != model.RoleApplicant {
		return nil, ErrForbidden
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return nil, err
	}

	status := model.ContactStatusPending
	opts := repository.ContantListOptions{
		RecipientID: &applicant.ID,
		Status:      &status,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
	}

	return s.repo.List(ctx, opts)
}

type UpdateStatusDTO struct {
	Auth   AuthContext
	ID     uuid.UUID
	Status model.ContactStatus
}

func (s *contactService) UpdateStatus(ctx context.Context, dto UpdateStatusDTO) error {
	if dto.Auth.Role != model.RoleApplicant {
		return ErrForbidden
	}
	if dto.Status != model.ContactStatusAccepted && dto.Status != model.ContactStatusRejected {
		return ErrInvalidInput
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	contact, err := s.repo.GetByID(ctx, dto.ID)
	if err != nil {
		return err
	}

	if contact.RecipientID != applicant.ID {
		return ErrForbidden
	}

	updates := map[string]any{
		"status": dto.Status,
	}

	return s.repo.Update(ctx, dto.ID, updates)
}

type DeleteContactDTO struct {
	Auth AuthContext
	ID   uuid.UUID
}

func (s *contactService) Delete(ctx context.Context, dto DeleteContactDTO) error {
	if dto.Auth.Role != model.RoleApplicant {
		return ErrForbidden
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return err
	}

	contact, err := s.repo.GetByID(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if contact.SenderID != applicant.ID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, dto.ID)
}
