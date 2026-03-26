package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/pkg/auth"
	"github.com/nakle1ka/Tramplin/internal/pkg/hash"
	"github.com/nakle1ka/Tramplin/internal/repository"
)

type AuthService interface {
	Register(ctx context.Context, dto CreateAccountDTO) (AuthResult, error)
	Login(ctx context.Context, email, password string) (AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (AuthResult, error)

	GetRefreshExpires() time.Duration
	GetAccessExpires() time.Duration
}

type CreateAccountDTO struct {
	Email    string
	Password string
	Role     model.Role

	Employer  *model.Employer
	Applicant *model.Applicant
}

type AuthResult struct {
	AccessToken  string
	RefreshToken string
	UserID       uuid.UUID
	Role         model.Role
	IsVerified   bool
}

type authService struct {
	userRepo      repository.UserRepository
	applicantRepo repository.ApplicantRepository
	curatorRepo   repository.CuratorRepository
	employerRepo  repository.EmployerRepository
	cacheRepo     repository.CacheRepository
	txManager     repository.TransactionManager

	passwordHasher hash.Hasher
	tokenHasher    hash.Hasher
	tokenManager   auth.TokenManager

	refreshExp time.Duration
	accessExp  time.Duration
}

func (s *authService) Register(ctx context.Context, dto CreateAccountDTO) (AuthResult, error) {
	if dto.Role == model.RoleEmployer {
		valid, err := validateEmployerINN(dto.Employer.INN)
		if !valid || err != nil {
			return AuthResult{}, ErrInvalidEmployerINN
		}
	}

	var response AuthResult

	err := s.txManager.Wrap(ctx, func(txCtx context.Context) error {
		hashedPassword, err := s.passwordHasher.Hash([]byte(dto.Password))
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		user := &model.User{
			Email:        dto.Email,
			PasswordHash: string(hashedPassword),
			Role:         dto.Role,
			IsVerified:   false,
		}

		if err := s.userRepo.Create(txCtx, user); err != nil {
			return err
		}

		if err := s.createProfile(txCtx, user.ID, dto); err != nil {
			return err
		}

		accessToken, refreshToken, err := s.issueSession(user.ID, user.Role)
		if err != nil {
			return err
		}

		response = AuthResult{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			UserID:       user.ID,
			Role:         user.Role,
			IsVerified:   user.IsVerified,
		}

		return nil
	})

	return response, err
}

func (s *authService) createProfile(ctx context.Context, userID uuid.UUID, dto CreateAccountDTO) error {
	switch dto.Role {
	case model.RoleApplicant:
		if dto.Applicant == nil {
			return fmt.Errorf("applicant profile is required")
		}

		dto.Applicant.UserID = userID
		return s.applicantRepo.Create(ctx, dto.Applicant)

	case model.RoleEmployer:
		if dto.Employer == nil {
			return fmt.Errorf("employer profile is required")
		}

		dto.Employer.UserID = userID
		dto.Employer.VerifiedStatus = model.StatusPending
		return s.employerRepo.Create(ctx, dto.Employer)

	default:
		return fmt.Errorf("unknown role: %v", dto.Role)
	}
}

func (s *authService) Login(ctx context.Context, email string, password string) (AuthResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	if !s.passwordHasher.Verify([]byte(password), []byte(user.PasswordHash)) {
		return AuthResult{}, ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.issueSession(user.ID, user.Role)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Role:         user.Role,
		IsVerified:   user.IsVerified,
	}, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return fmt.Errorf("validate token: %w", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	key := fmt.Sprintf("session:%v:%v", userID, claims.TokenId)
	return s.cacheRepo.Delete(key)
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	claims, err := s.tokenManager.ValidateToken(refreshToken)
	if err != nil {
		return AuthResult{}, fmt.Errorf("validate token: %w", err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return AuthResult{}, fmt.Errorf("parse user id: %w", err)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return AuthResult{}, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return AuthResult{}, ErrUserNotFound
	}

	key := fmt.Sprintf("session:%v:%v", userID, claims.TokenId)
	storedHash, err := s.cacheRepo.Get(key)
	if err != nil {
		return AuthResult{}, fmt.Errorf("get session: %w", err)
	}

	if !s.tokenHasher.Verify([]byte(refreshToken), []byte(storedHash)) {
		return AuthResult{}, ErrInvalidToken
	}

	var response AuthResult
	err = s.txManager.Wrap(ctx, func(txCtx context.Context) error {
		_ = s.cacheRepo.Delete(key)

		accessToken, newRefreshToken, err := s.issueSession(userID, user.Role)
		if err != nil {
			return err
		}

		response = AuthResult{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			UserID:       user.ID,
			Role:         user.Role,
			IsVerified:   user.IsVerified,
		}

		return nil
	})

	return response, err
}

func (s *authService) issueSession(userID uuid.UUID, userRole model.Role) (string, string, error) {
	accessToken, _, err := s.tokenManager.GenerateToken(auth.TokenDTO{
		UserID:   userID.String(),
		UserRole: userRole,
		Expires:  s.accessExp,
	})
	if err != nil {
		return "", "", fmt.Errorf("generate access: %w", err)
	}

	refreshToken, tokenID, err := s.tokenManager.GenerateToken(auth.TokenDTO{
		UserID:   userID.String(),
		UserRole: userRole,
		Expires:  s.refreshExp,
	})
	if err != nil {
		return "", "", fmt.Errorf("generate refresh: %w", err)
	}

	tokenHash, err := s.tokenHasher.Hash([]byte(refreshToken))
	if err != nil {
		return "", "", fmt.Errorf("hash refresh: %w", err)
	}
	hashStr := hex.EncodeToString(tokenHash)

	key := fmt.Sprintf("session:%v:%v", userID, tokenID)
	if err := s.cacheRepo.Set(key, hashStr, s.refreshExp); err != nil {
		return "", "", fmt.Errorf("save session: %w", err)
	}

	return accessToken, refreshToken, nil
}

func validateEmployerINN(inn string) (bool, error) {
	n := len(inn)
	if n != 10 && n != 12 {
		return false, errors.New("INN must be 10 or 12 digits")
	}

	digits := make([]int, n)
	for i := 0; i < n; i++ {
		if inn[i] < '0' || inn[i] > '9' {
			return false, errors.New("invalid character")
		}
		digits[i] = int(inn[i] - '0')
	}

	check := func(d []int, weights []int) int {
		sum := 0
		for i, w := range weights {
			sum += d[i] * w
		}
		return (sum % 11) % 10
	}

	if n == 10 {
		w10 := []int{2, 4, 10, 3, 5, 9, 4, 6, 8}
		return digits[9] == check(digits[:9], w10), nil
	} else {
		w11 := []int{7, 2, 4, 10, 3, 5, 9, 4, 6, 8}
		w12 := []int{3, 7, 2, 4, 10, 3, 5, 9, 4, 6, 8}

		c11 := check(digits[:10], w11)
		c12 := check(digits[:11], w12)

		return digits[10] == c11 && digits[11] == c12, nil
	}
}

func (s *authService) GetRefreshExpires() time.Duration {
	return s.refreshExp
}

func (s *authService) GetAccessExpires() time.Duration {
	return s.accessExp
}

type Option func(*authService)

func WithAccessExpires(exp int) Option {
	return func(s *authService) {
		s.accessExp = time.Duration(exp) * time.Second
	}
}

func WithRefreshExpires(exp int) Option {
	return func(s *authService) {
		s.refreshExp = time.Duration(exp) * time.Second
	}
}

func NewAuthService(
	ur repository.UserRepository,
	ar repository.ApplicantRepository,
	cr repository.CuratorRepository,
	er repository.EmployerRepository,
	cacheRepo repository.CacheRepository,
	tm repository.TransactionManager,
	th hash.Hasher,
	ph hash.Hasher,
	tk auth.TokenManager,
	opts ...Option,
) AuthService {
	srv := &authService{
		userRepo:       ur,
		applicantRepo:  ar,
		curatorRepo:    cr,
		employerRepo:   er,
		cacheRepo:      cacheRepo,
		txManager:      tm,
		tokenHasher:    th,
		passwordHasher: ph,
		tokenManager:   tk,
		refreshExp:     time.Hour * 24 * 7,
		accessExp:      time.Minute * 15,
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}
