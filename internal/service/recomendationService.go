package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type RecomendationService interface {
	Create(ctx context.Context, dto CreateRecomendationDTO) error
	GetMyRecomendations(ctx context.Context, dto ListRecomendationsDTO) ([]model.Recomendation, error)
	Delete(ctx context.Context, dto DeleteRecomendationDTO) error
}

type recomendationService struct {
	repo            repository.RecomendationRepository
	applicantRepo   repository.ApplicantRepository
	opportinityRepo repository.OpportunityRepository
	contactRepo     repository.ContactRepository
}

func NewRecomendationService(
	repo repository.RecomendationRepository,
	applicantRepo repository.ApplicantRepository,
	opportinityRepo repository.OpportunityRepository,
	contactRepo repository.ContactRepository,
) RecomendationService {
	return &recomendationService{
		repo:            repo,
		applicantRepo:   applicantRepo,
		opportinityRepo: opportinityRepo,
		contactRepo:     contactRepo,
	}
}

type CreateRecomendationDTO struct {
	Auth         AuthContext
	RecipientID  uuid.UUID
	OportunityID uuid.UUID
}

func (s *recomendationService) Create(ctx context.Context, dto CreateRecomendationDTO) error {
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

	_, contactErr := s.contactRepo.GetByApplicantsIDs(ctx, sender.ID, dto.RecipientID)
	if contactErr != nil {
		if errors.Is(contactErr, repository.ErrNotFound) {
			return ErrForbidden
		}
		return contactErr
	}

	opportinity, err := s.opportinityRepo.GetByID(ctx, dto.OportunityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrInvalidInput
		}
		return err
	}

	rec := &model.Recomendation{
		ID:            uuid.New(),
		SenderID:      sender.ID,
		RecipientID:   dto.RecipientID,
		OpportunityID: opportinity.ID,
	}

	return s.repo.Create(ctx, rec)
}

type ListRecomendationsDTO struct {
	Auth   AuthContext
	Limit  int
	Offset int
}

func (s *recomendationService) GetMyRecomendations(ctx context.Context, dto ListRecomendationsDTO) ([]model.Recomendation, error) {
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

	opts := repository.RecomendationListOptions{
		RecipientID: &applicant.ID,
		Limit:       dto.Limit,
		Offset:      dto.Offset,
	}

	return s.repo.List(ctx, opts)
}

type DeleteRecomendationDTO struct {
	Auth AuthContext
	ID   uuid.UUID
}

func (s *recomendationService) Delete(ctx context.Context, dto DeleteRecomendationDTO) error {
	if dto.Auth.Role != model.RoleApplicant {
		return ErrForbidden
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return err
	}

	recomendation, err := s.repo.GetByID(ctx, dto.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	if recomendation.SenderID != applicant.ID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, dto.ID)
}
