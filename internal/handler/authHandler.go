package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/service"
	"gorm.io/gorm"
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
	var authCtx *service.AuthContext
	extracted, err := extractAuthContext(c)
	if err == nil {
		authCtx = &extracted
	}

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil || !req.Role.IsValid() {
		slog.Error("failed to bind register request", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	createDTO, err := mapRegisterRequestToDTO(req)
	if err != nil {
		slog.Error("failed to map register request", "error", err, "email", req.Email)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	createDTO.Auth = authCtx

	resp, err := h.authService.Register(c.Request.Context(), createDTO)
	if err != nil {
		slog.Error("failed to register user", "error", err, "email", req.Email)
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "email already exists"})
		case errors.Is(err, service.ErrInvalidEmployerINN):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid employer INN"})
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		slog.Warn("missing credentials", "email", req.Email)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "email and password are required"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		slog.Error("failed to login", "error", err, "email", req.Email)
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("user logged in successfully", "user_id", resp.UserID, "role", resp.Role)

	h.setRefreshTokenCookie(c, resp.RefreshToken)

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
		c.JSON(http.StatusOK, dto.ErrorResponse{Error: "logged out successfully"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), service.LogoutRequest{
		RefreshToken: refreshToken,
	}); err != nil {
		slog.Error("failed to logout", "error", err)
	}

	h.clearRefreshTokenCookie(c)

	slog.Info("user logged out successfully")
	c.JSON(http.StatusOK, dto.ErrorResponse{Error: "logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		slog.Warn("no refresh token in cookies")
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "refresh token not found"})
		return
	}

	resp, err := h.authService.Refresh(c.Request.Context(), service.RefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		slog.Error("failed to refresh token", "error", err)
		h.clearRefreshTokenCookie(c)
		switch {
		case errors.Is(err, service.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid refresh token"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
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

func mapRegisterRequestToDTO(req dto.RegisterRequest) (service.RegisterRequest, error) {
	dto := service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	switch req.Role {
	case model.RoleApplicant:
		if req.Applicant == nil {
			return service.RegisterRequest{}, errors.New("applicant data is required")
		}
		dto.Applicant = &model.Applicant{
			FirstName:  req.Applicant.FirstName,
			SecondName: req.Applicant.SecondName,
			LastName:   req.Applicant.LastName,
		}
	case model.RoleEmployer:
		if req.Employer == nil {
			return service.RegisterRequest{}, errors.New("employer data is required")
		}
		dto.Employer = &model.Employer{
			CompanyName: req.Employer.CompanyName,
			Description: req.Employer.Description,
			INN:         req.Employer.INN,
			Website:     req.Employer.Website,
		}
	}

	return dto, nil
}
