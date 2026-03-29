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

type OpportunityHandler struct {
	opportunityService service.OpportunityService
}

func NewOpportunityHandler(opportunityService service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{
		opportunityService: opportunityService,
	}
}

func (h *OpportunityHandler) Create(c *gin.Context) {
	var req dto.CreateOpportunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind create opportunity request", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := service.CreateOpportunityRequest{
		Auth:            authCtx,
		TagNames:        req.TagNames,
		Title:           req.Title,
		Description:     req.Description,
		OpportunityType: req.OpportunityType,
		WorkFormat:      req.WorkFormat,
		LocationCity:    req.LocationCity,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		ExperienceLevel: req.ExperienceLevel,
		ExpiresAt:       req.ExpiresAt,
		EventDateStart:  req.EventDateStart,
		EventDateEnd:    req.EventDateEnd,
	}

	opportunity, err := h.opportunityService.Create(c.Request.Context(), serviceReq)
	if err != nil {
		slog.Error("failed to create opportunity", "error", err, "user_id", authCtx.UserID)
		switch {
		case errors.Is(err, service.ErrInvalidOpportunity):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity data"})
		case errors.Is(err, service.ErrInvalidSalary):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid salary: min must be less than or equal to max"})
		case errors.Is(err, service.ErrInvalidDates):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid dates: event start must be before event end"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("opportunity created", "id", opportunity.ID, "user_id", authCtx.UserID)
	c.JSON(http.StatusCreated, h.toOpportunityResponse(opportunity))
}

func (h *OpportunityHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid opportunity id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
		return
	}

	var req dto.UpdateOpportunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind update opportunity request", "error", err, "opportunity_id", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := service.UpdateOpportunityRequest{
		Auth:            authCtx,
		ID:              id,
		Title:           req.Title,
		Description:     req.Description,
		OpportunityType: req.OpportunityType,
		WorkFormat:      req.WorkFormat,
		LocationCity:    req.LocationCity,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		ExperienceLevel: req.ExperienceLevel,
		ExpiresAt:       req.ExpiresAt,
		EventDateStart:  req.EventDateStart,
		EventDateEnd:    req.EventDateEnd,
	}

	opportunity, err := h.opportunityService.Update(c.Request.Context(), serviceReq)
	if err != nil {
		slog.Error("failed to update opportunity", "error", err, "opportunity_id", id, "user_id", authCtx.UserID)
		switch {
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "opportunity not found"})
		case errors.Is(err, service.ErrExpiredOpportunity):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "cannot modify expired opportunity"})
		case errors.Is(err, service.ErrAlreadyModerated):
			c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "opportunity already moderated"})
		case errors.Is(err, service.ErrInvalidOpportunity):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity data"})
		case errors.Is(err, service.ErrInvalidSalary):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid salary: min must be less than or equal to max"})
		case errors.Is(err, service.ErrInvalidDates):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid dates: event start must be before event end"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("opportunity updated", "id", id, "user_id", authCtx.UserID)
	c.JSON(http.StatusOK, h.toOpportunityResponse(opportunity))
}

func (h *OpportunityHandler) UpdateModerationStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid opportunity id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
		return
	}

	var req dto.UpdateModerationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind moderation update request", "error", err, "opportunity_id", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := service.UpdateModerationStatusRequest{
		Auth:   authCtx,
		ID:     id,
		Status: req.Status,
	}

	if err := h.opportunityService.UpdateModerationStatus(c.Request.Context(), serviceReq); err != nil {
		slog.Error("failed to update moderation status", "error", err, "opportunity_id", id, "user_id", authCtx.UserID)
		switch {
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "only curators can update moderation status"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "opportunity not found"})
		case errors.Is(err, service.ErrInvalidInput):
			c.Status(http.StatusBadRequest)
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("moderation status updated", "id", id, "status", req.Status, "curator_id", authCtx.UserID)
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *OpportunityHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid opportunity id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
		return
	}

	var authCtx *service.AuthContext
	extractedAuthCtx, err := extractAuthContext(c)
	if err != nil {
		authCtx = nil
	} else {
		authCtx = &extractedAuthCtx
	}

	serviceReq := service.GetByIdOpportunityRequest{
		Auth: authCtx,
		ID:   id,
	}

	opportunity, err := h.opportunityService.GetByID(c.Request.Context(), serviceReq)
	if err != nil {
		slog.Error("failed to get opportunity by id", "error", err, "opportunity_id", id)
		switch {
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "opportunity not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, h.toOpportunityResponse(opportunity))
}

func (h *OpportunityHandler) List(c *gin.Context) {
	var req dto.ListOpportunitiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("failed to bind list query", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	var EmployerID, CuratorID *uuid.UUID
	if req.EmployerID != nil {
		parsed, err := uuid.Parse(*req.EmployerID)
		if err != nil {
			slog.Warn("invalid employer id format", "id", *req.EmployerID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid employer id"})
			return
		}
		EmployerID = &parsed
	}
	if req.CuratorID != nil {
		parsed, err := uuid.Parse(*req.CuratorID)
		if err != nil {
			slog.Warn("invalid curator id format", "id", *req.CuratorID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid curator id"})
			return
		}
		CuratorID = &parsed
	}

	var authCtx *service.AuthContext
	extractedAuthCtx, err := extractAuthContext(c)
	if err != nil {
		authCtx = nil
	} else {
		authCtx = &extractedAuthCtx
	}

	filter := service.OpportunityFilter{
		EmployerID:       EmployerID,
		CuratorID:        CuratorID,
		OpportunityType:  req.OpportunityType,
		WorkFormat:       req.WorkFormat,
		ExperienceLevel:  req.ExperienceLevel,
		ModerationStatus: req.ModerationStatus,
		LocationCity:     req.LocationCity,
		TagIDs:           req.TagIDs,
		MinSalary:        req.MinSalary,
		MaxSalary:        req.MaxSalary,
		ExpiresAfter:     req.ExpiresAfter,
		Limit:            req.Limit,
		Offset:           req.Offset,
	}

	serviceReq := service.ListRequest{
		Auth:   authCtx,
		Filter: filter,
	}

	opportunities, total, err := h.opportunityService.List(c.Request.Context(), serviceReq)
	if err != nil {
		slog.Error("failed to list opportunities", "error", err)
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	response := dto.ListOpportunitiesResponse{
		Opportunities: make([]dto.OpportunityResponse, len(opportunities)),
		Total:         total,
		Limit:         filter.Limit,
		Offset:        filter.Offset,
	}

	for i, opp := range opportunities {
		response.Opportunities[i] = h.toOpportunityResponse(opp)
	}

	c.JSON(http.StatusOK, response)
}

func (h *OpportunityHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid opportunity id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := service.DeleteRequest{
		Auth: authCtx,
		ID:   id,
	}

	if err := h.opportunityService.Delete(c.Request.Context(), serviceReq); err != nil {
		slog.Error("failed to delete opportunity", "error", err, "opportunity_id", id, "user_id", authCtx.UserID)
		switch {
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "opportunity not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("opportunity deleted", "id", id, "user_id", authCtx.UserID)
	c.JSON(http.StatusNoContent, nil)
}

func (h *OpportunityHandler) AddTags(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid opportunity id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
		return
	}

	var req dto.AddTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind add tags request", "error", err, "opportunity_id", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := service.AddTagsRequest{
		Auth:          authCtx,
		OpportunityID: id,
		TagIDs:        req.TagIDs,
	}

	if err := h.opportunityService.AddTags(c.Request.Context(), serviceReq); err != nil {
		slog.Error("failed to add tags to opportunity", "error", err, "opportunity_id", id)
		switch {
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "opportunity not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("tags added to opportunity", "id", id, "tag_ids", req.TagIDs, "user_id", authCtx.UserID)
	c.JSON(http.StatusOK, gin.H{"status": "tags added"})
}

func (h *OpportunityHandler) RemoveTags(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.Warn("invalid opportunity id format", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
		return
	}

	var req dto.RemoveTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind remove tags request", "error", err, "opportunity_id", id)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	serviceReq := service.RemoveTagsRequest{
		Auth:          authCtx,
		OpportunityID: id,
		TagIDs:        req.TagIDs,
	}

	if err := h.opportunityService.RemoveTags(c.Request.Context(), serviceReq); err != nil {
		slog.Error("failed to remove tags from opportunity", "error", err, "opportunity_id", id)
		switch {
		case errors.Is(err, service.ErrForbidden):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
		case errors.Is(err, service.ErrNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "opportunity not found"})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("tags removed from opportunity", "id", id, "tag_ids", req.TagIDs, "user_id", authCtx.UserID)
	c.JSON(http.StatusOK, gin.H{"status": "tags removed"})
}

func (h *OpportunityHandler) toOpportunityResponse(opp *model.Opportunity) dto.OpportunityResponse {
	tags := make([]dto.TagResponse, len(opp.Tags))
	for i, tag := range opp.Tags {
		tags[i] = dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return dto.OpportunityResponse{
		ID:               opp.ID,
		EmployerID:       opp.EmployerID,
		CuratorID:        opp.CuratorID,
		Tags:             tags,
		Title:            opp.Title,
		Description:      opp.Description,
		OpportunityType:  opp.OpportunityType,
		WorkFormat:       opp.WorkFormat,
		LocationCity:     opp.LocationCity,
		Latitude:         opp.Latitude,
		Longitude:        opp.Longitude,
		SalaryMin:        opp.SalaryMin,
		SalaryMax:        opp.SalaryMax,
		ExperienceLevel:  opp.ExperienceLevel,
		ExpiresAt:        opp.ExpiresAt,
		EventDateStart:   opp.EventDateStart,
		EventDateEnd:     opp.EventDateEnd,
		ModerationStatus: opp.ModerationStatus,
		CreatedAt:        opp.CreatedAt,
		UpdatedAt:        opp.UpdatedAt,
	}
}
