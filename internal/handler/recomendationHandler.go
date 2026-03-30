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

type RecommendationHandler struct {
	service service.RecomendationService
}

func NewRecommendationHandler(s service.RecomendationService) *RecommendationHandler {
	return &RecommendationHandler{service: s}
}

func (h *RecommendationHandler) Create(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("failed to extract auth context", "err", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreateRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid create recommendation body", "err", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	err = h.service.Create(c.Request.Context(), service.CreateRecomendationDTO{
		Auth:         authCtx,
		RecipientID:  req.RecipientID,
		OportunityID: req.OpportunityID,
	})

	if err != nil {
		slog.Error("service failed to create recommendation", "err", err, "user_id", authCtx.UserID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *RecommendationHandler) GetMyRecomendations(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("failed to extract auth context", "err", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var query dto.ListRecommendationsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Warn("invalid query params", "err", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query parameters"})
		return
	}

	recs, err := h.service.GetMyRecomendations(c.Request.Context(), service.ListRecomendationsDTO{
		Auth:   authCtx,
		Limit:  query.Limit,
		Offset: query.Offset,
	})
	if err != nil {
		slog.Error("service failed to get recommendations", "err", err, "user_id", authCtx.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	response := make([]dto.RecomendationResponse, len(recs))
	for i, rec := range recs {
		response[i] = dto.RecomendationResponse{
			ID: rec.ID,

			SenderID: rec.SenderID,
			Sender: dto.RecomendationInfo{
				FirstName:  rec.Sender.FirstName,
				SecondName: rec.Sender.SecondName,
				LastName:   rec.Sender.LastName,
			},

			RecipientID: rec.RecipientID,
			Recipient: dto.RecomendationInfo{
				FirstName:  rec.Recipient.FirstName,
				SecondName: rec.Recipient.SecondName,
				LastName:   rec.Recipient.LastName,
			},

			OpportunityID: rec.OpportunityID,
			Opportunity: dto.OpportunityRecomendationInfo{
				Title:       rec.Opportunity.Title,
				Description: rec.Opportunity.Description,
			},
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *RecommendationHandler) Delete(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("failed to extract auth context", "err", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid uuid in path", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id format"})
		return
	}

	err = h.service.Delete(c.Request.Context(), service.DeleteRecomendationDTO{
		Auth: authCtx,
		ID:   id,
	})

	if err != nil {
		slog.Error("service failed to delete recommendation", "err", err, "rec_id", id)
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "recommendation not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "forbidden"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
