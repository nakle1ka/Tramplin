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

type OpportunityService interface {
	Create(ctx context.Context, req CreateOpportunityRequest) (*model.Opportunity, error)
	Update(ctx context.Context, req UpdateOpportunityRequest) (*model.Opportunity, error)
	UpdateModerationStatus(ctx context.Context, req UpdateModerationStatusRequest) error
	Delete(ctx context.Context, req DeleteRequest) error

	List(ctx context.Context, req ListRequest) ([]*model.Opportunity, int64, error)
	GetByID(ctx context.Context, req GetByIdOpportunityRequest) (*model.Opportunity, error)
	GetByEmployerID(ctx context.Context, req GetByEmployerIDRequest) ([]*model.Opportunity, error)
	GetByCuratorID(ctx context.Context, req GetByCuratorIDRequest) ([]*model.Opportunity, error)

	AddTags(ctx context.Context, req AddTagsRequest) error
	RemoveTags(ctx context.Context, req RemoveTagsRequest) error
}

type opportunityService struct {
	opportunityRepo repository.OpportunityRepository
	tagRepo         repository.TagRepository
	employerRepo    repository.EmployerRepository
	curatorRepo     repository.CuratorRepository
}

func NewOpportunityService(
	opportunityRepo repository.OpportunityRepository,
	tagRepo repository.TagRepository,
	employerRepo repository.EmployerRepository,
	curatorRepo repository.CuratorRepository,
) OpportunityService {
	return &opportunityService{
		opportunityRepo: opportunityRepo,
		tagRepo:         tagRepo,
		employerRepo:    employerRepo,
		curatorRepo:     curatorRepo,
	}
}

type CreateOpportunityRequest struct {
	Auth            AuthContext
	TagNames        []string
	Title           string
	Description     string
	OpportunityType model.OpportunityType
	WorkFormat      model.WorkFormat
	LocationCity    string
	Latitude        float64
	Longitude       float64
	SalaryMin       int
	SalaryMax       int
	ExperienceLevel model.ExpirienseLevel
	ExpiresAt       *time.Time
	EventDateStart  *time.Time
	EventDateEnd    *time.Time
}

func (s *opportunityService) Create(ctx context.Context, req CreateOpportunityRequest) (*model.Opportunity, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, ErrInvalidInput
	}

	var employerID, curatorID *uuid.UUID
	switch req.Auth.Role {
	case model.RoleEmployer:
		employer, err := s.employerRepo.GetByUserID(ctx, req.Auth.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get employer: %w", err)
		}
		employerID = &employer.ID

	case model.RoleCurator:
		curator, err := s.curatorRepo.GetByUserID(ctx, req.Auth.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to get curator: %w", err)
		}
		curatorID = &curator.ID

	default:
		return nil, ErrForbidden
	}

	var tags []*model.Tag
	if len(req.TagNames) > 0 {
		var err error
		tags, err = s.tagRepo.GetOrCreateBatch(ctx, req.TagNames)
		if err != nil {
			return nil, fmt.Errorf("failed to process tags: %w", err)
		}
	}

	opportunity := &model.Opportunity{
		ID:               uuid.New(),
		EmployerID:       employerID,
		CuratorID:        curatorID,
		Tags:             tags,
		Title:            req.Title,
		Description:      req.Description,
		OpportunityType:  req.OpportunityType,
		WorkFormat:       req.WorkFormat,
		LocationCity:     req.LocationCity,
		Latitude:         req.Latitude,
		Longitude:        req.Longitude,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		ExperienceLevel:  req.ExperienceLevel,
		ExpiresAt:        req.ExpiresAt,
		EventDateStart:   req.EventDateStart,
		EventDateEnd:     req.EventDateEnd,
		ModerationStatus: model.ModerationStatusPending,
	}

	if err := s.opportunityRepo.Create(ctx, opportunity); err != nil {
		return nil, fmt.Errorf("failed to create opportunity: %w", err)
	}

	return opportunity, nil
}

type UpdateOpportunityRequest struct {
	Auth            AuthContext
	ID              uuid.UUID
	Title           *string
	Description     *string
	OpportunityType *model.OpportunityType
	WorkFormat      *model.WorkFormat
	LocationCity    *string
	Latitude        *float64
	Longitude       *float64
	SalaryMin       *int
	SalaryMax       *int
	ExperienceLevel *model.ExpirienseLevel
	ExpiresAt       *time.Time
	EventDateStart  *time.Time
	EventDateEnd    *time.Time
}

func (s *opportunityService) Update(ctx context.Context, req UpdateOpportunityRequest) (*model.Opportunity, error) {
	existing, err := s.opportunityRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if err := s.checkOwnership(existing, req.Auth); err != nil {
		return nil, err
	}

	if existing.ExpiresAt != nil && existing.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpiredOpportunity
	}

	if existing.ModerationStatus != model.ModerationStatusPending {
		return nil, ErrAlreadyModerated
	}

	updates := make(map[string]any)

	if req.Title != nil {
		if err := s.validateTitle(*req.Title); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidOpportunity, err)
		}
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		if err := s.validateDescription(*req.Description); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidOpportunity, err)
		}
		updates["description"] = *req.Description
	}
	if req.OpportunityType != nil {
		updates["opportunity_type"] = *req.OpportunityType
	}
	if req.WorkFormat != nil {
		updates["work_format"] = *req.WorkFormat
	}
	if req.LocationCity != nil {
		updates["location_city"] = *req.LocationCity
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.SalaryMin != nil || req.SalaryMax != nil {
		salaryMin := existing.SalaryMin
		salaryMax := existing.SalaryMax

		if req.SalaryMin != nil {
			salaryMin = *req.SalaryMin
			updates["salary_min"] = salaryMin
		}
		if req.SalaryMax != nil {
			salaryMax = *req.SalaryMax
			updates["salary_max"] = salaryMax
		}

		if err := s.validateSalary(salaryMin, salaryMax); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidOpportunity, err)
		}
	}
	if req.ExperienceLevel != nil {
		updates["experience_level"] = *req.ExperienceLevel
	}
	if req.ExpiresAt != nil {
		updates["expires_at"] = *req.ExpiresAt
	}
	if req.EventDateStart != nil || req.EventDateEnd != nil {
		eventStart := existing.EventDateStart
		eventEnd := existing.EventDateEnd

		if req.EventDateStart != nil {
			eventStart = req.EventDateStart
			updates["event_date_start"] = eventStart
		}
		if req.EventDateEnd != nil {
			eventEnd = req.EventDateEnd
			updates["event_date_end"] = eventEnd
		}

		if err := s.validateEventDates(eventStart, eventEnd); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidOpportunity, err)
		}
	}

	if len(updates) == 0 {
		return existing, nil
	}

	if err := s.opportunityRepo.Update(ctx, req.ID, updates); err != nil {
		return nil, err
	}

	return s.opportunityRepo.GetByID(ctx, req.ID)
}

type UpdateModerationStatusRequest struct {
	Auth   AuthContext
	ID     uuid.UUID
	Status model.ModerationStatus
}

func (s *opportunityService) UpdateModerationStatus(ctx context.Context, req UpdateModerationStatusRequest) error {
	if req.Auth.Role != model.RoleCurator {
		return ErrForbidden
	}

	if !req.Status.IsValid() {
		return ErrInvalidInput
	}

	return s.opportunityRepo.Update(ctx, req.ID, map[string]any{
		"moderation_status": req.Status,
	})
}

type GetByIdOpportunityRequest struct {
	Auth *AuthContext
	ID   uuid.UUID
}

func (s *opportunityService) GetByID(ctx context.Context, req GetByIdOpportunityRequest) (*model.Opportunity, error) {
	opportunity, err := s.opportunityRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if opportunity.ModerationStatus != model.ModerationStatusApproved {
		if req.Auth == nil {
			return nil, ErrForbidden
		}
		if err := s.checkOwnership(opportunity, *req.Auth); err != nil {
			return nil, err
		}
	}

	return opportunity, nil
}

type GetByEmployerIDRequest struct {
	Auth       AuthContext
	EmployerID uuid.UUID
}

func (s *opportunityService) GetByEmployerID(ctx context.Context, req GetByEmployerIDRequest) ([]*model.Opportunity, error) {
	if req.Auth.Role != model.RoleCurator && req.Auth.UserID != req.EmployerID {
		return nil, ErrForbidden
	}

	return s.opportunityRepo.GetByEmployerID(ctx, req.EmployerID)
}

type GetByCuratorIDRequest struct {
	Auth      AuthContext
	CuratorID uuid.UUID
}

func (s *opportunityService) GetByCuratorID(ctx context.Context, req GetByCuratorIDRequest) ([]*model.Opportunity, error) {
	if req.Auth.Role != model.RoleCurator && req.Auth.UserID != req.CuratorID {
		return nil, ErrForbidden
	}

	return s.opportunityRepo.GetByCuratorID(ctx, req.CuratorID)
}

type ListRequest struct {
	Auth   *AuthContext
	Filter OpportunityFilter
}

type OpportunityFilter struct {
	EmployerID       *uuid.UUID
	CuratorID        *uuid.UUID
	OpportunityType  *model.OpportunityType
	WorkFormat       *model.WorkFormat
	ExperienceLevel  *model.ExpirienseLevel
	ModerationStatus *model.ModerationStatus
	LocationCity     *string
	TagIDs           []uuid.UUID
	MinSalary        *int
	MaxSalary        *int
	ExpiresAfter     *bool
	Limit            int
	Offset           int
}

func (s *opportunityService) List(ctx context.Context, req ListRequest) ([]*model.Opportunity, int64, error) {
	if req.Filter.Limit <= 0 || req.Filter.Limit > 100 || req.Filter.Offset < 0 {
		return nil, 0, ErrInvalidInput
	}

	if req.Auth == nil || req.Auth.Role != model.RoleCurator {
		approved := model.ModerationStatusApproved
		req.Filter.ModerationStatus = &approved
	}

	repoFilter := repository.OpportunityFilter{
		EmployerID:       req.Filter.EmployerID,
		CuratorID:        req.Filter.CuratorID,
		OpportunityType:  req.Filter.OpportunityType,
		WorkFormat:       req.Filter.WorkFormat,
		ExperienceLevel:  req.Filter.ExperienceLevel,
		ModerationStatus: req.Filter.ModerationStatus,
		LocationCity:     req.Filter.LocationCity,
		TagIDs:           req.Filter.TagIDs,
		MinSalary:        req.Filter.MinSalary,
		MaxSalary:        req.Filter.MaxSalary,
		ExpiresAfter:     req.Filter.ExpiresAfter,
		Limit:            req.Filter.Limit,
		Offset:           req.Filter.Offset,
	}

	return s.opportunityRepo.List(ctx, repoFilter)
}

type DeleteRequest struct {
	Auth AuthContext
	ID   uuid.UUID
}

func (s *opportunityService) Delete(ctx context.Context, req DeleteRequest) error {
	existing, err := s.opportunityRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if err := s.checkOwnership(existing, req.Auth); err != nil {
		return err
	}

	return s.opportunityRepo.Delete(ctx, req.ID)
}

type AddTagsRequest struct {
	Auth          AuthContext
	OpportunityID uuid.UUID
	TagIDs        []uuid.UUID
}

func (s *opportunityService) AddTags(ctx context.Context, req AddTagsRequest) error {
	opportunity, err := s.opportunityRepo.GetByID(ctx, req.OpportunityID)
	if err != nil {
		return err
	}

	if err := s.checkOwnership(opportunity, req.Auth); err != nil {
		return err
	}

	tags, err := s.tagRepo.GetByIDs(ctx, req.TagIDs)
	if err != nil {
		return fmt.Errorf("failed to get tags: %w", err)
	}

	if len(tags) != len(req.TagIDs) {
		return errors.New("some tags not found")
	}

	return s.opportunityRepo.AddTags(ctx, req.OpportunityID, req.TagIDs)
}

type RemoveTagsRequest struct {
	Auth          AuthContext
	OpportunityID uuid.UUID
	TagIDs        []uuid.UUID
}

func (s *opportunityService) RemoveTags(ctx context.Context, req RemoveTagsRequest) error {
	opportunity, err := s.opportunityRepo.GetByID(ctx, req.OpportunityID)
	if err != nil {
		return err
	}

	if err := s.checkOwnership(opportunity, req.Auth); err != nil {
		return err
	}

	return s.opportunityRepo.RemoveTags(ctx, req.OpportunityID, req.TagIDs)
}

func (s *opportunityService) checkOwnership(opportunity *model.Opportunity, auth AuthContext) error {
	if auth.Role == model.RoleCurator {
		return nil
	}

	if opportunity.EmployerID != nil && *opportunity.EmployerID == auth.UserID {
		return nil
	}

	if opportunity.CuratorID != nil && *opportunity.CuratorID == auth.UserID {
		return nil
	}

	return ErrForbidden
}

func (s *opportunityService) validateCreateRequest(req CreateOpportunityRequest) error {
	if req.Title == "" {
		return errors.New("title is required")
	}
	if err := s.validateTitle(req.Title); err != nil {
		return err
	}

	if req.Description == "" {
		return errors.New("description is required")
	}
	if err := s.validateDescription(req.Description); err != nil {
		return err
	}

	if err := s.validateSalary(req.SalaryMin, req.SalaryMax); err != nil {
		return err
	}

	if err := s.validateEventDates(req.EventDateStart, req.EventDateEnd); err != nil {
		return err
	}

	return nil
}

func (s *opportunityService) validateTitle(title string) error {
	if len(title) > 60 {
		return errors.New("title must not exceed 60 characters")
	}
	if len(title) < 5 {
		return errors.New("title must be at least 3 characters")
	}
	return nil
}

func (s *opportunityService) validateDescription(description string) error {
	if len(description) < 10 {
		return errors.New("description must be at least 10 characters")
	}
	return nil
}

func (s *opportunityService) validateSalary(min, max int) error {
	if min < 0 || max < 0 {
		return errors.New("salary values cannot be negative")
	}
	if min > max {
		return ErrInvalidSalary
	}
	return nil
}

func (s *opportunityService) validateEventDates(start, end *time.Time) error {
	if start != nil && end != nil && start.After(*end) {
		return ErrInvalidDates
	}
	return nil
}
