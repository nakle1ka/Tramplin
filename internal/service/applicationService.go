package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type ApplicationService interface {
	CreateApplication(ctx context.Context, dto CreateApplicationDTO) (*model.Application, error)
	UpdateApplicationStatus(ctx context.Context, dto UpdateApplicationStatusDTO) error
	DeleteApplication(ctx context.Context, dto DeleteApplicationDTO) error
	GetApplicationByID(ctx context.Context, dto GetApplicationByIDDTO) (*model.Application, error)
	GetApplicationsByOpportunity(ctx context.Context, dto GetApplicationsByOpportunityDTO) ([]*model.Application, int64, error)
	GetApplicationsByApplicant(ctx context.Context, dto GetApplicationsByApplicantDTO) ([]*model.Application, int64, error)
}

type applicationService struct {
	applicationRepo repository.ApplicationRepository
	opportnityRepo  repository.OpportunityRepository
	applicantRepo   repository.ApplicantRepository
	employerRepo    repository.EmployerRepository
}

func NewApplicationService(
	appRepo repository.ApplicationRepository,
	oppRepo repository.OpportunityRepository,
	applicRepo repository.ApplicantRepository,
	empRepo repository.EmployerRepository,
) ApplicationService {
	return &applicationService{
		applicationRepo: appRepo,
		opportnityRepo:  oppRepo,
		applicantRepo:   applicRepo,
		employerRepo:    empRepo,
	}
}

type CreateApplicationDTO struct {
	Auth          AuthContext
	OpportunityID uuid.UUID
}

func (s *applicationService) CreateApplication(ctx context.Context, dto CreateApplicationDTO) (*model.Application, error) {
	if dto.Auth.Role != model.RoleApplicant {
		return nil, ErrForbidden
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return nil, err
	}

	opportunity, err := s.opportnityRepo.GetByID(ctx, dto.OpportunityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if opportunity.ExpiresAt != nil && opportunity.ExpiresAt.Before(time.Now()) {
		return nil, ErrOpportunityClosed
	}
	if opportunity.EventDateEnd != nil && opportunity.EventDateEnd.Before(time.Now()) {
		return nil, ErrOpportunityClosed
	}

	application := &model.Application{
		ID:            uuid.New(),
		ApplicantID:   applicant.ID,
		OpportunityID: dto.OpportunityID,
		Status:        model.AcceptStatusPending,
	}

	if err := s.applicationRepo.Create(ctx, application); err != nil {
		return nil, err
	}

	return application, nil
}

type UpdateApplicationStatusDTO struct {
	Auth          AuthContext
	ApplicationID uuid.UUID
	Status        model.AcceptStatus
}

func (s *applicationService) UpdateApplicationStatus(ctx context.Context, dto UpdateApplicationStatusDTO) error {
	if dto.Auth.Role != model.RoleEmployer && dto.Auth.Role != model.RoleCurator {
		return ErrForbidden
	}

	application, err := s.applicationRepo.GetByID(ctx, dto.ApplicationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	switch dto.Auth.Role {
	case model.RoleEmployer:
		employer, err := s.employerRepo.GetByUserID(ctx, dto.Auth.UserID)
		if err != nil {
			return err
		}
		if application.Opportunity.EmployerID == nil || *application.Opportunity.EmployerID != employer.ID {
			return ErrForbidden
		}

	case model.RoleCurator:
		// кураторам можно принимать отклики только с
		// возможностей созданными другими кураторами
		if application.Opportunity.CuratorID == nil {
			return ErrForbidden
		}
	}

	return s.applicationRepo.UpdateAcceptStatus(ctx, dto.ApplicationID, dto.Status)
}

type DeleteApplicationDTO struct {
	Auth          AuthContext
	ApplicationID uuid.UUID
}

func (s *applicationService) DeleteApplication(ctx context.Context, dto DeleteApplicationDTO) error {
	if dto.Auth.Role != model.RoleApplicant {
		return ErrForbidden
	}

	application, err := s.applicationRepo.GetByID(ctx, dto.ApplicationID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
	if err != nil {
		return err
	}

	if applicant.ID != application.ApplicantID {
		return ErrForbidden
	}

	return s.applicationRepo.Delete(ctx, dto.ApplicationID)
}

type GetApplicationByIDDTO struct {
	Auth          AuthContext
	ApplicationID uuid.UUID
}

func (s *applicationService) GetApplicationByID(ctx context.Context, dto GetApplicationByIDDTO) (*model.Application, error) {
	application, err := s.applicationRepo.GetByID(ctx, dto.ApplicationID)
	if err != nil {
		return nil, ErrNotFound
	}

	switch dto.Auth.Role {
	case model.RoleApplicant:
		applicant, err := s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
		if err != nil {
			return nil, err
		}
		if applicant.ID != application.ApplicantID {
			return nil, ErrForbidden
		}

	case model.RoleEmployer:
		employer, err := s.employerRepo.GetByUserID(ctx, dto.Auth.UserID)
		if err != nil {
			return nil, err
		}
		if application.Opportunity.EmployerID == nil || *application.Opportunity.EmployerID != employer.ID {
			return nil, ErrForbidden
		}
	}

	return application, nil
}

type GetApplicationsByOpportunityDTO struct {
	Auth          AuthContext
	OpportunityID uuid.UUID
	Limit         int
	Offset        int
}

func (s *applicationService) GetApplicationsByOpportunity(ctx context.Context, dto GetApplicationsByOpportunityDTO) ([]*model.Application, int64, error) {
	if dto.Auth.Role != model.RoleEmployer && dto.Auth.Role != model.RoleCurator {
		return nil, 0, ErrForbidden
	}

	opportunity, err := s.opportnityRepo.GetByID(ctx, dto.OpportunityID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, 0, ErrNotFound
		}
		return nil, 0, err
	}

	if dto.Auth.Role == model.RoleEmployer {
		employer, err := s.employerRepo.GetByUserID(ctx, dto.Auth.UserID)
		if err != nil {
			return nil, 0, err
		}
		if opportunity.EmployerID == nil || *opportunity.EmployerID != employer.ID {
			return nil, 0, ErrForbidden
		}
	}

	return s.applicationRepo.GetByOpportunityID(ctx, dto.OpportunityID, dto.Limit, dto.Offset)
}

type GetApplicationsByApplicantDTO struct {
	Auth        AuthContext
	ApplicantID uuid.UUID
	Limit       int
	Offset      int
}

func (s *applicationService) GetApplicationsByApplicant(ctx context.Context, dto GetApplicationsByApplicantDTO) ([]*model.Application, int64, error) {
	applicant, err := s.applicantRepo.GetByID(ctx, dto.ApplicantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, 0, ErrNotFound
		}
		return nil, 0, err
	}

	if applicant.PrivacySetting != model.PrivacyPublic {
		var reqApplicant *model.Applicant
		if dto.Auth.Role == model.RoleApplicant {
			reqApplicant, err = s.applicantRepo.GetByUserID(ctx, dto.Auth.UserID)
			if err != nil {
				return nil, 0, err
			}
		}

		if reqApplicant == nil || reqApplicant.ID != applicant.ID {
			return nil, 0, ErrForbidden
		}
	}

	return s.applicationRepo.GetByApplicantID(ctx, dto.ApplicantID, dto.Limit, dto.Offset)
}
