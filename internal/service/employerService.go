package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type EmployerService interface {
	List(ctx context.Context, req ListEmployersRequest) (*ListEmployersResponse, error)
	GetMe(ctx context.Context, req GetMeEmployerRequest) (*model.Employer, error)
	GetByID(ctx context.Context, req GetByIdOpportunityRequest) (*model.Employer, error)
	Update(ctx context.Context, req UpdateEmployerRequest) error
	UpdateVerifiedStatus(ctx context.Context, req UpdateVerifiedStatusRequest) error
}

type employerService struct {
	repo repository.EmployerRepository
}

func NewEmployerService(repo repository.EmployerRepository) EmployerService {
	return &employerService{repo: repo}
}

type GetMeEmployerRequest struct {
	Auth AuthContext
}

func (s *employerService) GetMe(ctx context.Context, req GetMeEmployerRequest) (*model.Employer, error) {
	employer, err := s.repo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		return nil, err
	}
	return employer, nil
}

type GetByIDEmployerRequest struct {
	ID uuid.UUID
}

func (s *employerService) GetByID(ctx context.Context, req GetByIdOpportunityRequest) (*model.Employer, error) {
	employer, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return employer, nil
}

type UpdateEmployerRequest struct {
	Auth    AuthContext
	Request UpdateEmployerFields
}

type UpdateEmployerFields struct {
	CompanyName *string `json:"company_name"`
	Description *string `json:"description"`
	Website     *string `json:"website"`
}

func (s *employerService) Update(ctx context.Context, req UpdateEmployerRequest) error {
	employer, err := s.repo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		return err
	}

	if req.Auth.Role != model.RoleCurator && req.Auth.UserID != employer.UserID {
		return ErrForbidden
	}

	updates := make(map[string]any)

	fields := req.Request
	if fields.CompanyName != nil {
		updates["company_name"] = *fields.CompanyName
	}
	if fields.Description != nil {
		updates["description"] = *fields.Description
	}
	if fields.Website != nil {
		updates["website"] = *fields.Website
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.Update(ctx, employer.ID, updates)
}

type UpdateVerifiedStatusRequest struct {
	Auth   AuthContext
	ID     uuid.UUID
	Status model.VerificationStatus
}

func (s *employerService) UpdateVerifiedStatus(ctx context.Context, req UpdateVerifiedStatusRequest) error {
	if req.Auth.Role != model.RoleCurator {
		return ErrForbidden
	}

	employer, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	field := make(map[string]any)
	field["verified_status"] = req.Status

	return s.repo.Update(ctx, employer.ID, field)
}

type ListEmployersRequest struct {
	Filters ListEmployersFilters
	Auth    AuthContext
}

type ListEmployersFilters struct {
	CompanyName    *string
	VerifiedStatus *model.VerificationStatus
	Limit          int
	Offset         int
}

type ListEmployersResponse struct {
	Employers []*model.Employer
	Total     int64
	Limit     int
	Offset    int
}

func (s *employerService) List(ctx context.Context, req ListEmployersRequest) (*ListEmployersResponse, error) {
	if req.Auth.Role != model.RoleCurator {
		return nil, ErrForbidden
	}

	repoFilters := repository.EmployerListFilters{
		CompanyName:    req.Filters.CompanyName,
		VerifiedStatus: req.Filters.VerifiedStatus,
		Limit:          req.Filters.Limit,
		Offset:         req.Filters.Offset,
	}

	employers, total, err := s.repo.List(ctx, repoFilters)
	if err != nil {
		return nil, err
	}

	return &ListEmployersResponse{
		Employers: employers,
		Total:     total,
		Limit:     req.Filters.Limit,
		Offset:    req.Filters.Offset,
	}, nil
}
