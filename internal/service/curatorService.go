package service

import (
	"context"
	"errors"

	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type CuratorService struct {
	curatorRepo repository.CuratorRepository
}

func NewCuratorService(curatorRepo repository.CuratorRepository) *CuratorService {
	return &CuratorService{
		curatorRepo: curatorRepo,
	}
}

type GetMeRequest struct {
	Auth AuthContext
}

func (s *CuratorService) GetMe(ctx context.Context, req GetMeRequest) (*model.Curator, error) {
	if req.Auth.Role != model.RoleCurator {
		return nil, ErrForbidden
	}

	curator, err := s.curatorRepo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	if curator == nil {
		return nil, ErrNotFound
	}

	return curator, nil
}

type UpdateRequest struct {
	Auth     AuthContext
	FullName *string
}

func (s *CuratorService) Update(ctx context.Context, req UpdateRequest) error {
	curator, err := s.curatorRepo.GetByUserID(ctx, req.Auth.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	updates := make(map[string]any)

	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}

	if len(updates) == 0 {
		return nil
	}

	err = s.curatorRepo.Update(ctx, curator.ID, updates)
	if err != nil {
		return err
	}

	return nil
}
