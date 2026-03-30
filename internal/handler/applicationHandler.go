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

type ApplicationHandler struct {
	appService service.ApplicationService
}

func NewApplicationHandler(appService service.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{
		appService: appService,
	}
}

func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("invalid request body for create application", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for create application", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	app, err := h.appService.CreateApplication(c.Request.Context(), service.CreateApplicationDTO{
		Auth:          auth,
		OpportunityID: req.OpportunityID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			slog.Warn("forbidden to create application", "error", err, "user_id", auth.UserID, "opportunity_id", req.OpportunityID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrNotFound):
			slog.Warn("opportunity not found for application", "error", err, "user_id", auth.UserID, "opportunity_id", req.OpportunityID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrOpportunityClosed):
			slog.Warn("attempt to apply to closed opportunity", "user_id", auth.UserID, "opportunity_id", req.OpportunityID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		default:
			slog.Error("failed to create application", "error", err, "user_id", auth.UserID, "opportunity_id", req.OpportunityID)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("application created successfully", "application_id", app.ID, "user_id", auth.UserID, "opportunity_id", req.OpportunityID)
	c.JSON(http.StatusCreated, dto.ApplicationResponse{
		ID:            app.ID,
		OpportunityID: app.OpportunityID,
		ApplicantID:   app.ApplicantID,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,
	})
}

func (h *ApplicationHandler) UpdateApplicationStatus(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		slog.Warn("invalid application id for update status", "id", appIDStr, "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid application id"})
		return
	}

	var req dto.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || !req.Status.IsValid() {
		slog.Warn("invalid request body for update application status", "error", err, "status", req.Status)
		c.Status(http.StatusBadRequest)
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for update application status", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = h.appService.UpdateApplicationStatus(c.Request.Context(), service.UpdateApplicationStatusDTO{
		Auth:          auth,
		ApplicationID: appID,
		Status:        req.Status,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			slog.Warn("forbidden to update application status", "user_id", auth.UserID, "application_id", appID, "status", req.Status)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrNotFound):
			slog.Warn("application not found for status update", "application_id", appID, "user_id", auth.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			slog.Error("failed to update application status", "error", err, "user_id", auth.UserID, "application_id", appID, "status", req.Status)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("application status updated successfully", "application_id", appID, "user_id", auth.UserID, "status", req.Status)
	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		slog.Warn("invalid application id for delete", "id", appIDStr, "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid application id"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for delete application", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	err = h.appService.DeleteApplication(c.Request.Context(), service.DeleteApplicationDTO{
		Auth:          auth,
		ApplicationID: appID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			slog.Warn("forbidden to delete application", "user_id", auth.UserID, "application_id", appID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrNotFound):
			slog.Warn("application not found for delete", "application_id", appID, "user_id", auth.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			slog.Error("failed to delete application", "error", err, "user_id", auth.UserID, "application_id", appID)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("application deleted successfully", "application_id", appID, "user_id", auth.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}

func (h *ApplicationHandler) GetApplicationByID(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		slog.Warn("invalid application id for get by id", "id", appIDStr, "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid application id"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for get application", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	app, err := h.appService.GetApplicationByID(c.Request.Context(), service.GetApplicationByIDDTO{
		Auth:          auth,
		ApplicationID: appID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrForbidden):
			slog.Warn("forbidden to get application", "user_id", auth.UserID, "application_id", appID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrNotFound):
			slog.Warn("application not found", "application_id", appID, "user_id", auth.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			slog.Error("failed to get application by id", "error", err, "user_id", auth.UserID, "application_id", appID)
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	slog.Info("application retrieved successfully", "application_id", appID, "user_id", auth.UserID)
	c.JSON(http.StatusOK, dto.ApplicationResponse{
		ID:            app.ID,
		OpportunityID: app.OpportunityID,
		ApplicantID:   app.ApplicantID,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,
		Opportunity: &dto.OpportunityApplicationInfo{
			ID:           app.OpportunityID,
			Title:        app.Opportunity.Title,
			Description:  app.Opportunity.Description,
			LocationCity: app.Opportunity.LocationCity,
		},
		Applicant: &dto.ApplicationApplicantInfo{
			ID:        app.Applicant.ID,
			FirstName: app.Applicant.FirstName,
			LastName:  app.Applicant.LastName,
		},
	})
}

func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	var req dto.GetApplicationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("invalid query parameters for get applications", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for get applications", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if req.OpportunityID != nil {
		id, err := uuid.Parse(*req.OpportunityID)
		if err != nil {
			slog.Warn("invalid opportunity id in query", "opportunity_id", *req.OpportunityID, "error", err)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
			return
		}

		applications, total, err := h.appService.GetApplicationsByOpportunity(c.Request.Context(), service.GetApplicationsByOpportunityDTO{
			Auth:          auth,
			OpportunityID: id,
			Limit:         req.Limit,
			Offset:        req.Offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, service.ErrForbidden):
				slog.Warn("forbidden to get applications by opportunity", "user_id", auth.UserID, "opportunity_id", id)
				c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			case errors.Is(err, service.ErrNotFound):
				slog.Warn("opportunity not found for applications", "opportunity_id", id, "user_id", auth.UserID)
				c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			default:
				slog.Error("failed to get applications by opportunity", "error", err, "user_id", auth.UserID, "opportunity_id", id)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
			}
			return
		}

		slog.Info("applications retrieved by opportunity", "user_id", auth.UserID, "opportunity_id", id, "count", len(applications), "total", total)

		resp := make([]*dto.ApplicationResponse, len(applications))
		for i, app := range applications {
			resp[i] = &dto.ApplicationResponse{
				ID:            app.ID,
				OpportunityID: app.OpportunityID,
				ApplicantID:   app.ApplicantID,
				Status:        app.Status,
				CreatedAt:     app.CreatedAt,
				UpdatedAt:     app.UpdatedAt,
				Opportunity: &dto.OpportunityApplicationInfo{
					ID:           app.OpportunityID,
					Title:        app.Opportunity.Title,
					Description:  app.Opportunity.Description,
					LocationCity: app.Opportunity.LocationCity,
				},
				Applicant: &dto.ApplicationApplicantInfo{
					ID:        app.Applicant.ID,
					FirstName: app.Applicant.FirstName,
					LastName:  app.Applicant.LastName,
				},
			}
		}

		c.JSON(http.StatusOK, dto.ApplicationsListResponse{
			Applications: resp,
			Total:        total,
			Limit:        req.Limit,
			Offset:       req.Offset,
		})
		return
	}

	if req.ApplicantID != nil {
		id, err := uuid.Parse(*req.ApplicantID)
		if err != nil {
			slog.Warn("invalid applicant id in query", "applicant_id", *req.ApplicantID, "error", err)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
			return
		}

		applications, total, err := h.appService.GetApplicationsByApplicant(c.Request.Context(), service.GetApplicationsByApplicantDTO{
			Auth:        auth,
			ApplicantID: id,
			Limit:       req.Limit,
			Offset:      req.Offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, service.ErrForbidden):
				slog.Warn("forbidden to get applications by applicant", "user_id", auth.UserID, "applicant_id", id)
				c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			case errors.Is(err, service.ErrNotFound):
				slog.Warn("applicant not found for applications", "applicant_id", id, "user_id", auth.UserID)
				c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			default:
				slog.Error("failed to get applications by applicant", "error", err, "user_id", auth.UserID, "applicant_id", id)
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
			}
			return
		}

		slog.Info("applications retrieved by applicant", "user_id", auth.UserID, "applicant_id", id, "count", len(applications), "total", total)

		resp := make([]*dto.ApplicationResponse, len(applications))
		for i, app := range applications {
			resp[i] = &dto.ApplicationResponse{
				ID:            app.ID,
				OpportunityID: app.OpportunityID,
				ApplicantID:   app.ApplicantID,
				Status:        app.Status,
				CreatedAt:     app.CreatedAt,
				UpdatedAt:     app.UpdatedAt,
				Opportunity: &dto.OpportunityApplicationInfo{
					ID:           app.OpportunityID,
					Title:        app.Opportunity.Title,
					Description:  app.Opportunity.Description,
					LocationCity: app.Opportunity.LocationCity,
				},
				Applicant: &dto.ApplicationApplicantInfo{
					ID:        app.Applicant.ID,
					FirstName: app.Applicant.FirstName,
					LastName:  app.Applicant.LastName,
				},
			}
		}

		c.JSON(http.StatusOK, dto.ApplicationsListResponse{
			Applications: resp,
			Total:        total,
			Limit:        req.Limit,
			Offset:       req.Offset,
		})
		return
	}

	slog.Warn("missing required query parameter", "user_id", auth.UserID)
	c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "either opportunity_id or applicant_id is required"})
}
