package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type FavoritesHandler struct {
	service service.FavoritesService
}

func NewFavoritesHandler(svc service.FavoritesService) *FavoritesHandler {
	return &FavoritesHandler{service: svc}
}

func (h *FavoritesHandler) Create(c *gin.Context) {
	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("unauthorized create attempt", "err", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreateFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	fav, err := h.service.Create(c.Request.Context(), service.CreateFavoritesRequest{
		Auth:          auth,
		OpportunityID: req.OpportunityID,
	})

	if err != nil {
		slog.Error("failed to create favorite", "err", err, "user_id", auth.UserID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	slog.Info("favorite created", "id", fav.ID, "user_id", auth.UserID)
	c.JSON(http.StatusCreated, MapFavoriteToDTO(fav))
}

func (h *FavoritesHandler) List(c *gin.Context) {
	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("unauthorized list attempt", "err", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.ListFavoritesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("invalid query params", "err", err, "user_id", auth.UserID)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query parameters"})
		return
	}

	var applicantID, opportunityID *uuid.UUID
	if req.ApplicantID != nil {
		id, err := uuid.Parse(*req.ApplicantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id format"})
			return
		}
		applicantID = &id
	}
	if req.OpportunityID != nil {
		id, err := uuid.Parse(*req.OpportunityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id format"})
			return
		}
		opportunityID = &id
	}

	favs, total, err := h.service.List(c.Request.Context(), service.ListFavoritesRequest{
		Auth:          auth,
		ApplicantID:   applicantID,
		OpportunityID: opportunityID,
		Limit:         req.Limit,
		Offset:        req.Offset,
	})

	if err != nil {
		slog.Error("failed to list favorites", "err", err, "user_id", auth.UserID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	items := make([]dto.FavoriteResponse, 0, len(favs))
	for _, f := range favs {
		items = append(items, MapFavoriteToDTO(&f))
	}

	c.JSON(http.StatusOK, dto.FavoriteListResponse{
		Items: items,
		Total: total,
	})
}

func (h *FavoritesHandler) Delete(c *gin.Context) {
	auth, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id format"})
		return
	}

	err = h.service.Delete(c.Request.Context(), service.DeleteFavoritesRequest{
		Auth: auth,
		ID:   id,
	})

	if err != nil {
		slog.Error("failed to delete favorite", "err", err, "id", id)
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal error"})
		return
	}

	slog.Info("favorite deleted", "id", id, "user_id", auth.UserID)
	c.Status(http.StatusNoContent)
}

func MapFavoriteToDTO(m *model.Favorites) dto.FavoriteResponse {
	return dto.FavoriteResponse{
		ID:            m.ID,
		ApplicantID:   m.ApplicantID,
		OpportunityID: m.OpportunityID,
		Opportunity: dto.OpportunityFavoritesResponse{
			Title:       m.Opportunity.Title,
			Description: m.Opportunity.Description,
			SalaryMin:   m.Opportunity.SalaryMin,
			SalaryMax:   m.Opportunity.SalaryMax,
		},
	}
}
