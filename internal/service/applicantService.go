package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type ApplicantService interface {
	GetMe(ctx context.Context, req GetMeApplicantRequest) (*model.Applicant, error)
	GetByID(ctx context.Context, req GetApplicantByIDRequest) (*model.Applicant, error)

	Update(ctx context.Context, req UpdateApplicantRequest) error

	GetTags(ctx context.Context, req GetApplicantTagsRequest) ([]*model.Tag, error)
	SetTags(ctx context.Context, req SetApplicantTagsRequest) error
	AddTags(ctx context.Context, req AddApplicantTagsRequest) error
	RemoveTags(ctx context.Context, req RemoveApplicantTagsRequest) error

	List(ctx context.Context, req ListApplicantsRequest) ([]*model.Applicant, error)
}

type applicantService struct {
	repo repository.ApplicantRepository
}

func NewApplicantService(repo repository.ApplicantRepository) ApplicantService {
	return &applicantService{
		repo: repo,
	}
}

type GetMeApplicantRequest struct {
	Auth AuthContext
}

func (s *applicantService) GetMe(ctx context.Context, req GetMeApplicantRequest) (*model.Applicant, error) {
	applicant, err := s.repo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get applicant: %w", err)
	}

	return applicant, nil
}

type ListApplicantsRequest struct {
	Limit  int
	Offset int
}

func (s *applicantService) List(ctx context.Context, req ListApplicantsRequest) ([]*model.Applicant, error) {
	limit := req.Limit
	offset := req.Offset

	if limit <= 0 || limit > 100 || offset < 0 {
		return nil, ErrInvalidInput
	}

	applicants, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list applicants: %w", err)
	}

	return applicants, nil
}

type GetApplicantByIDRequest struct {
	Auth *AuthContext
	ID   uuid.UUID
}

func (s *applicantService) GetByID(ctx context.Context, req GetApplicantByIDRequest) (*model.Applicant, error) {
	applicant, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get applicant: %w", err)
	}

	if applicant.PrivacySetting == model.PrivacyPrivate && (req.Auth == nil || req.Auth.UserID != applicant.UserID) {
		applicant = &model.Applicant{
			ID:         applicant.ID,
			UserID:     applicant.UserID,
			User:       applicant.User,
			FirstName:  applicant.FirstName,
			SecondName: applicant.SecondName,
			LastName:   applicant.LastName,
		}
	}

	return applicant, nil
}

type UpdateApplicantRequest struct {
	Auth           AuthContext
	FirstName      *string
	SecondName     *string
	LastName       *string
	University     *string
	GraduationYear *int
	About          *string
	PrivacySetting *model.Privacy
	WorkExpirience *string
}

func (s *applicantService) Update(ctx context.Context, req UpdateApplicantRequest) error {
	applicant, err := s.repo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get applicant: %w", err)
	}

	if req.Auth.Role != model.RoleCurator && applicant.UserID != req.Auth.UserID {
		return ErrForbidden
	}

	updates := make(map[string]any)

	if req.FirstName != nil {
		if err := s.validateName(*req.FirstName, "first name"); err != nil {
			return err
		}
		updates["first_name"] = *req.FirstName
	}
	if req.SecondName != nil {
		if err := s.validateName(*req.SecondName, "second name"); err != nil {
			return err
		}
		updates["second_name"] = *req.SecondName
	}
	if req.LastName != nil {
		if err := s.validateName(*req.LastName, "last name"); err != nil {
			return err
		}
		updates["last_name"] = *req.LastName
	}
	if req.University != nil {
		if err := s.validateUniversity(*req.University); err != nil {
			return err
		}
		updates["university"] = *req.University
	}
	if req.GraduationYear != nil {
		if err := s.validateGraduationYear(*req.GraduationYear); err != nil {
			return err
		}
		updates["graduation_year"] = *req.GraduationYear
	}
	if req.About != nil {
		if err := s.validateAbout(*req.About); err != nil {
			return err
		}
		updates["about"] = *req.About
	}
	if req.PrivacySetting != nil {
		if !req.PrivacySetting.IsValid() {
			return ErrInvalidInput
		}
		updates["privacy_setting"] = *req.PrivacySetting
	}
	if req.WorkExpirience != nil {
		updates["work_experience"] = *req.WorkExpirience
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.repo.Update(ctx, applicant.ID, updates); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update applicant: %w", err)
	}

	return nil
}

type GetApplicantTagsRequest struct {
	Auth        *AuthContext
	ApplicantID uuid.UUID
}

func (s *applicantService) GetTags(ctx context.Context, req GetApplicantTagsRequest) ([]*model.Tag, error) {
	applicant, err := s.repo.GetByID(ctx, req.ApplicantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get applicant: %w", err)
	}

	if applicant.PrivacySetting == model.PrivacyPrivate {
		if req.Auth == nil {
			return nil, ErrForbidden
		}
		if req.Auth.UserID != applicant.UserID && req.Auth.Role != model.RoleCurator {
			return nil, ErrForbidden
		}
	}

	tags, err := s.repo.GetTags(ctx, req.ApplicantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}

type SetApplicantTagsRequest struct {
	Auth        AuthContext
	ApplicantID uuid.UUID
	TagIDs      []uuid.UUID
}

func (s *applicantService) SetTags(ctx context.Context, req SetApplicantTagsRequest) error {
	applicant, err := s.repo.GetByID(ctx, req.ApplicantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get applicant: %w", err)
	}

	if applicant.UserID != req.Auth.UserID {
		return ErrForbidden
	}

	if err := s.repo.SetTags(ctx, req.ApplicantID, req.TagIDs); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to set tags: %w", err)
	}

	return nil
}

type AddApplicantTagsRequest struct {
	Auth        AuthContext
	ApplicantID uuid.UUID
	TagIDs      []uuid.UUID
}

func (s *applicantService) AddTags(ctx context.Context, req AddApplicantTagsRequest) error {
	if len(req.TagIDs) == 0 {
		return nil
	}

	applicant, err := s.repo.GetByID(ctx, req.ApplicantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get applicant: %w", err)
	}

	if applicant.UserID != req.Auth.UserID {
		return ErrForbidden
	}

	if err := s.repo.AddTags(ctx, req.ApplicantID, req.TagIDs); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to add tags: %w", err)
	}

	return nil
}

type RemoveApplicantTagsRequest struct {
	Auth        AuthContext
	ApplicantID uuid.UUID
	TagIDs      []uuid.UUID
}

func (s *applicantService) RemoveTags(ctx context.Context, req RemoveApplicantTagsRequest) error {
	if len(req.TagIDs) == 0 {
		return nil
	}

	applicant, err := s.repo.GetByID(ctx, req.ApplicantID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to get applicant: %w", err)
	}

	if applicant.UserID != req.Auth.UserID {
		return ErrForbidden
	}

	if err := s.repo.RemoveTags(ctx, req.ApplicantID, req.TagIDs); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to remove tags: %w", err)
	}

	return nil
}

func (s *applicantService) validateName(name, field string) error {
	if len(name) > 100 {
		return fmt.Errorf("%s must not exceed 100 characters", field)
	}
	if len(name) < 1 {
		return fmt.Errorf("%s cannot be empty", field)
	}
	return nil
}

func (s *applicantService) validateUniversity(university string) error {
	if len(university) > 200 {
		return errors.New("university name must not exceed 200 characters")
	}
	return nil
}

func (s *applicantService) validateGraduationYear(year int) error {
	currentYear := time.Now().Year()
	if year < 1950 || year > currentYear+10 {
		return errors.New("invalid graduation year")
	}
	return nil
}

func (s *applicantService) validateAbout(about string) error {
	if len(about) > 5000 {
		return errors.New("about text must not exceed 5000 characters")
	}
	return nil
}
