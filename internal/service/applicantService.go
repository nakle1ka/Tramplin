package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type Applicant = model.Applicant

type ApplicantService interface {
	GetMe(ctx context.Context, userId uuid.UUID) (*Applicant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Applicant, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateApplicantRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	AddTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error
	RemoveTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error
	GetTags(ctx context.Context, applicantID uuid.UUID) ([]*model.Tag, error)
	SetTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error
}

type applicantService struct {
	repo repository.ApplicantRepository
}

func (s *applicantService) GetMe(ctx context.Context, userId uuid.UUID) (*Applicant, error) {
	applicant, err := s.repo.GetByUserID(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return nil, ErrApplicantNotFound
		}
		return nil, fmt.Errorf("failed to get applicant: %w", err)
	}

	return applicant, nil
}

func NewApplicantService(repo repository.ApplicantRepository) ApplicantService {
	return &applicantService{
		repo: repo,
	}
}

func (s *applicantService) GetByID(ctx context.Context, id uuid.UUID) (*Applicant, error) {
	applicant, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return nil, ErrApplicantNotFound
		}
		return nil, fmt.Errorf("failed to get applicant: %w", err)
	}

	return applicant, nil
}

func (s *applicantService) Update(ctx context.Context, id uuid.UUID, req UpdateApplicantRequest) error {

	updates := make(map[string]any)

	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.SecondName != nil {
		updates["second_name"] = *req.SecondName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.University != nil {
		updates["university"] = *req.University
	}
	if req.GraduationYear != nil {
		updates["graduation_year"] = *req.GraduationYear
	}
	if req.About != nil {
		updates["about"] = *req.About
	}
	if req.PrivacySetting != nil {
		if req.PrivacySetting.IsValid() {
			updates["privacy_setting"] = *req.PrivacySetting
		} else {
			return ErrInvalidInput
		}
	}

	if len(updates) == 0 {
		return nil
	}

	if err := s.repo.Update(ctx, id, updates); err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return ErrApplicantNotFound
		}
		return fmt.Errorf("failed to update applicant: %w", err)
	}

	return nil
}

func (s *applicantService) Delete(ctx context.Context, id uuid.UUID) error {

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return ErrApplicantNotFound
		}
		return fmt.Errorf("failed to delete applicant: %w", err)
	}

	return nil
}

func (s *applicantService) AddTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}

	if err := s.repo.AddTags(ctx, applicantID, tagIDs); err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return ErrApplicantNotFound
		}
		return fmt.Errorf("failed to add tags: %w", err)
	}

	return nil
}

func (s *applicantService) RemoveTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error {
	if len(tagIDs) == 0 {
		return nil
	}

	if err := s.repo.RemoveTags(ctx, applicantID, tagIDs); err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return ErrApplicantNotFound
		}
		return fmt.Errorf("failed to remove tags: %w", err)
	}

	return nil
}

func (s *applicantService) GetTags(ctx context.Context, applicantID uuid.UUID) ([]*model.Tag, error) {

	tags, err := s.repo.GetTags(ctx, applicantID)
	if err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return nil, ErrApplicantNotFound
		}
		return nil, fmt.Errorf("failed to get tags: %w", err)
	}

	return tags, nil
}

func (s *applicantService) SetTags(ctx context.Context, applicantID uuid.UUID, tagIDs []uuid.UUID) error {

	if err := s.repo.SetTags(ctx, applicantID, tagIDs); err != nil {
		if errors.Is(err, repository.ErrApplicantNotFound) {
			return ErrApplicantNotFound
		}
		return fmt.Errorf("failed to set tags: %w", err)
	}

	return nil
}

type UpdateApplicantRequest struct {
	FirstName      *string        `json:"first_name,omitempty"`
	SecondName     *string        `json:"second_name,omitempty"`
	LastName       *string        `json:"last_name,omitempty"`
	University     *string        `json:"university,omitempty"`
	GraduationYear *int           `json:"graduation_year,omitempty"`
	About          *string        `json:"about,omitempty"`
	PrivacySetting *model.Privacy `json:"privacy_setting,omitempty"`
}
