package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type FavoritesService interface {
	Create(ctx context.Context, req CreateFavoritesRequest) (*model.Favorites, error)
	GetByID(ctx context.Context, req GetFavoritesRequest) (*model.Favorites, error)
	List(ctx context.Context, req ListFavoritesRequest) ([]model.Favorites, int64, error)
	Delete(ctx context.Context, req DeleteFavoritesRequest) error
}

type favoritesService struct {
	favoritesRepo   repository.FavoritesRepository
	applicantRepo   repository.ApplicantRepository
	opportunityRepo repository.OpportunityRepository
}

func NewFavoritesService(
	favoritesRepo repository.FavoritesRepository,
	applicantRepo repository.ApplicantRepository,
	opportunityRepo repository.OpportunityRepository,
) FavoritesService {
	return &favoritesService{
		favoritesRepo:   favoritesRepo,
		applicantRepo:   applicantRepo,
		opportunityRepo: opportunityRepo,
	}
}

type CreateFavoritesRequest struct {
	Auth          AuthContext
	OpportunityID uuid.UUID
}

func (s *favoritesService) Create(ctx context.Context, req CreateFavoritesRequest) (*model.Favorites, error) {
	if req.Auth.Role != model.RoleApplicant {
		return nil, ErrForbidden
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		return nil, err
	}

	opportunity, err := s.opportunityRepo.GetByID(ctx, req.OpportunityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	favorite := &model.Favorites{
		ID:            uuid.New(),
		ApplicantID:   applicant.ID,
		OpportunityID: opportunity.ID,
	}

	if err := s.favoritesRepo.Create(ctx, favorite); err != nil {
		return nil, err
	}

	return favorite, nil
}

type GetFavoritesRequest struct {
	auth AuthContext
	ID   uuid.UUID
}

func (s *favoritesService) GetByID(ctx context.Context, req GetFavoritesRequest) (*model.Favorites, error) {
	if req.auth.Role != model.RoleApplicant {
		return nil, ErrForbidden
	}

	favorite, err := s.favoritesRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, req.auth.UserID)
	if err != nil {
		return nil, err
	}
	if applicant.ID != favorite.ApplicantID {
		return nil, ErrForbidden
	}

	return favorite, nil
}

type ListFavoritesRequest struct {
	Auth          AuthContext
	ApplicantID   *uuid.UUID
	OpportunityID *uuid.UUID
	Limit         int
	Offset        int
}

func (s *favoritesService) List(ctx context.Context, req ListFavoritesRequest) ([]model.Favorites, int64, error) {
	if req.Auth.Role != model.RoleApplicant {
		return nil, 0, ErrForbidden
	}

	if req.ApplicantID != nil && req.Auth.Role != model.RoleCurator {
		applicant, err := s.applicantRepo.GetByUserID(ctx, req.Auth.UserID)
		if err != nil {
			return nil, 0, err
		}
		if applicant.ID != *req.ApplicantID {
			return nil, 0, ErrForbidden
		}
	} else if req.Auth.Role != model.RoleCurator {
		return nil, 0, ErrForbidden
	}

	opts := repository.FavoritesListOptions{
		ApplicantID:   req.ApplicantID,
		OpportunityID: req.OpportunityID,
		Limit:         req.Limit,
		Offset:        req.Offset,
	}

	return s.favoritesRepo.List(ctx, opts)
}

type DeleteFavoritesRequest struct {
	Auth AuthContext
	ID   uuid.UUID
}

func (s *favoritesService) Delete(ctx context.Context, req DeleteFavoritesRequest) error {
	if req.Auth.Role != model.RoleApplicant {
		return ErrForbidden
	}

	favorite, err := s.favoritesRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	applicant, err := s.applicantRepo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		return err
	}

	if applicant.ID != favorite.ApplicantID {
		return ErrForbidden
	}

	return s.favoritesRepo.Delete(ctx, req.ID)
}
