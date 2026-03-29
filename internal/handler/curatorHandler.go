package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type CuratorHandler struct {
	curatorService *service.CuratorService
}

func NewCuratorHandler(curatorService *service.CuratorService) *CuratorHandler {
	return &CuratorHandler{
		curatorService: curatorService,
	}
}

func (h *CuratorHandler) GetMe(c *gin.Context) {
	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	req := service.GetMeRequest{
		Auth: auth,
	}

	curator, err := h.curatorService.GetMe(c.Request.Context(), req)
	if err != nil {
		slog.Error("failed to get curator profile", "error", err, "user_id", auth.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.CuratorResponse{
		ID:           curator.ID,
		Email:        curator.User.Email,
		FullName:     curator.FullName,
		IsSuperAdmin: curator.IsSuperAdmin,
		CreatedAt:    curator.CreatedAt,
		UpdatedAt:    curator.UpdatedAt,
	})
}

func (h *CuratorHandler) Update(c *gin.Context) {
	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind curator update request", "error", err, "user_id", auth.UserID)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	serviceReq := service.UpdateRequest{
		Auth:     auth,
		FullName: req.FullName,
	}

	updErr := h.curatorService.Update(c.Request.Context(), serviceReq)
	if updErr != nil {
		slog.Error("failed to update curator profile", "error", updErr, "user_id", auth.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: updErr.Error()})
		return
	}

	slog.Info("curator profile updated", "user_id", auth.UserID)
	c.Status(http.StatusOK)
}
