package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type EmployerHandler struct {
	service service.EmployerService
}

func NewEmployerHandler(service service.EmployerService) *EmployerHandler {
	return &EmployerHandler{service: service}
}

func (h *EmployerHandler) GetMe(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context",
			slog.String("error", err.Error()),
		)
		c.Status(http.StatusUnauthorized)
		return
	}

	res, err := h.service.GetMe(c.Request.Context(), service.GetMeEmployerRequest{Auth: authCtx})
	if err != nil {
		slog.Error("failed to get employer info",
			slog.String("error", err.Error()),
			slog.String("user_id", authCtx.UserID.String()),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	response := dto.EmployerResponse{
		ID:             res.ID,
		UserID:         res.UserID,
		Email:          res.User.Email,
		CompanyName:    res.CompanyName,
		INN:            res.INN,
		Description:    res.Description,
		Website:        res.Website,
		VerifiedStatus: res.VerifiedStatus,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *EmployerHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid uuid parameter",
			slog.String("id", c.Param("id")),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid uuid",
		})
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), service.GetByIdOpportunityRequest{ID: id})
	if err != nil {
		slog.Info("employer not found",
			slog.String("employer_id", id.String()),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "employer not found",
		})
		return
	}

	response := dto.EmployerResponse{
		ID:             res.ID,
		UserID:         res.UserID,
		Email:          res.User.Email,
		CompanyName:    res.CompanyName,
		INN:            res.INN,
		Description:    res.Description,
		Website:        res.Website,
		VerifiedStatus: res.VerifiedStatus,
		CreatedAt:      res.CreatedAt,
		UpdatedAt:      res.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *EmployerHandler) Update(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context",
			slog.String("error", err.Error()),
		)
		c.Status(http.StatusUnauthorized)
		return
	}

	var req dto.UpdateEmployerFields
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid request body",
			slog.String("error", err.Error()),
			slog.String("user_id", authCtx.UserID.String()),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	err = h.service.Update(c.Request.Context(), service.UpdateEmployerRequest{
		Auth: authCtx,
		Request: service.UpdateEmployerFields{
			CompanyName: req.CompanyName,
			Description: req.Description,
			Website:     req.Website,
		},
	})

	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied for employer update",
				slog.String("user_id", authCtx.UserID.String()),
			)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "access denied",
			})
			return
		}
		slog.Error("failed to update employer",
			slog.String("error", err.Error()),
			slog.String("user_id", authCtx.UserID.String()),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	slog.Info("employer updated successfully",
		slog.String("user_id", authCtx.UserID.String()),
	)
	c.Status(http.StatusNoContent)
}

func (h *EmployerHandler) UpdateVerificationStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("invalid uuid parameter",
			slog.String("id", idStr),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid uuid",
		})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context",
			slog.String("error", err.Error()),
		)
		c.Status(http.StatusUnauthorized)
		return
	}

	var req dto.UpdateVerificationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid request body for verification status update",
			slog.String("error", err.Error()),
			slog.String("user_id", authCtx.UserID.String()),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}

	err = h.service.UpdateVerifiedStatus(c.Request.Context(), service.UpdateVerifiedStatusRequest{
		Auth:   authCtx,
		ID:     id,
		Status: *req.Status,
	})

	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("curator role required for verification status update",
				slog.String("user_id", authCtx.UserID.String()),
			)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "curator role required",
			})
			return
		}
		slog.Error("failed to update verification status",
			slog.String("error", err.Error()),
			slog.String("id", idStr),
			slog.Int("status", int(*req.Status)),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	slog.Info("verification status updated successfully",
		slog.String("user_id", authCtx.UserID.String()),
		slog.Int("status", int(*req.Status)),
	)
	c.Status(http.StatusOK)
}

func (h *EmployerHandler) List(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context",
			slog.String("error", err.Error()),
		)
		c.Status(http.StatusUnauthorized)
		return
	}

	var req dto.ListEmployersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("invalid query parameters",
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "invalid query parameters",
		})
		return
	}

	res, err := h.service.List(c.Request.Context(), service.ListEmployersRequest{
		Filters: service.ListEmployersFilters{
			CompanyName:    req.CompanyName,
			VerifiedStatus: req.VerifiedStatus,
			Limit:          req.Limit,
			Offset:         req.Offset,
		},
		Auth: authCtx,
	})

	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied for employers list",
				slog.String("user_id", authCtx.UserID.String()),
				slog.Int("role", int(authCtx.Role)),
			)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: "access denied",
			})
			return
		}
		slog.Error("failed to list employers",
			slog.String("error", err.Error()),
			slog.String("user_id", authCtx.UserID.String()),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "internal server error",
		})
		return
	}

	employersResponse := make([]dto.EmployerResponse, len(res.Employers))
	for i, employer := range res.Employers {
		employersResponse[i] = dto.EmployerResponse{
			ID:             employer.ID,
			UserID:         employer.UserID,
			Email:          employer.User.Email,
			CompanyName:    employer.CompanyName,
			INN:            employer.INN,
			Description:    employer.Description,
			Website:        employer.Website,
			VerifiedStatus: employer.VerifiedStatus,
			CreatedAt:      employer.CreatedAt,
			UpdatedAt:      employer.UpdatedAt,
		}
	}

	response := dto.ListEmployersResponse{
		Employers: employersResponse,
		Total:     res.Total,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	c.JSON(http.StatusOK, response)
}
