package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind register request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateRegisterRequest(req); err != nil {
		slog.Warn("invalid register request", "error", err, "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createDTO, err := mapRegisterRequestToDTO(req)
	if err != nil {
		slog.Error("failed to map register request", "error", err, "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), createDTO)
	if err != nil {
		slog.Error("failed to register user", "error", err, "email", req.Email)
		switch {
		case errors.Is(err, service.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		case errors.Is(err, service.ErrInvalidEmployerINN):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employer INN"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	slog.Info("user registered successfully", "user_id", resp.UserID, "role", resp.Role)

	h.setRefreshTokenCookie(c, resp.RefreshToken)

	c.JSON(http.StatusCreated, dto.AuthResponse{
		AccessToken: resp.AccessToken,
		UserID:      resp.UserID,
		Role:        resp.Role,
		IsVerified:  resp.IsVerified,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("failed to bind login request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		slog.Warn("missing credentials", "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("failed to login", "error", err, "email", req.Email)
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	slog.Info("user logged in successfully", "user_id", resp.UserID, "role", resp.Role)

	// Сохраняем только refresh token в cookie
	h.setRefreshTokenCookie(c, resp.RefreshToken)

	// Отдаем access token в ответе
	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken: resp.AccessToken,
		UserID:      resp.UserID,
		Role:        resp.Role,
		IsVerified:  resp.IsVerified,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("no refresh token in cookies")
		c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), refreshToken); err != nil {
		slog.Error("failed to logout", "error", err)
		// Даже если произошла ошибка, очищаем куки
	}

	h.clearRefreshTokenCookie(c)

	slog.Info("user logged out successfully")
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("no refresh token in cookies")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	resp, err := h.authService.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		slog.Error("failed to refresh token", "error", err)
		h.clearRefreshTokenCookie(c)
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	slog.Info("token refreshed successfully", "user_id", resp.UserID)

	h.setRefreshTokenCookie(c, resp.RefreshToken)

	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken: resp.AccessToken,
		UserID:      resp.UserID,
		Role:        resp.Role,
		IsVerified:  resp.IsVerified,
	})
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, refreshToken string) {
	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(h.authService.GetRefreshExpires().Seconds()),
		"/",
		"",
		true,
		true,
	)
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
}

func validateRegisterRequest(req dto.RegisterRequest) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if !req.Role.IsValid() {
		return errors.New("invalid role")
	}
	return nil
}

func mapRegisterRequestToDTO(req dto.RegisterRequest) (service.CreateAccountDTO, error) {
	dto := service.CreateAccountDTO{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	switch req.Role {
	case model.RoleApplicant:
		if req.Applicant == nil {
			return service.CreateAccountDTO{}, errors.New("applicant data is required")
		}
		dto.Applicant = &model.Applicant{
			FirstName:  req.Applicant.FirstName,
			SecondName: req.Applicant.SecondName,
			LastName:   req.Applicant.LastName,
		}
	case model.RoleEmployer:
		if req.Employer == nil {
			return service.CreateAccountDTO{}, errors.New("employer data is required")
		}
		dto.Employer = &model.Employer{
			CompanyName: req.Employer.CompanyName,
			Description: req.Employer.Description,
			Website:     req.Employer.Website,
		}
	}

	return dto, nil
}
