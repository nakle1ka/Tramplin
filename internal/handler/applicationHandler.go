package handler

import (
	"errors"
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
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
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrOpportunityNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrOpportunityClosed):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid application id"})
		return
	}

	var req dto.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || !req.Status.IsValid() {
		c.Status(http.StatusBadRequest)
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
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
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrApplicationNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid application id"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
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
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrApplicationNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}

func (h *ApplicationHandler) GetApplicationByID(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid application id"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
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
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
		case errors.Is(err, service.ErrApplicationNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ApplicationResponse{
		ID:            app.ID,
		OpportunityID: app.OpportunityID,
		ApplicantID:   app.ApplicantID,
		Status:        app.Status,
		CreatedAt:     app.CreatedAt,
		UpdatedAt:     app.UpdatedAt,
		Opportunity: &dto.OpportunityInfo{
			ID:           app.OpportunityID,
			Title:        app.Opportunity.Title,
			Description:  app.Opportunity.Description,
			LocationCity: app.Opportunity.LocationCity,
		},
		Applicant: &dto.ApplicantInfo{
			ID:        app.Applicant.ID,
			FirstName: app.Applicant.FirstName,
			LastName:  app.Applicant.LastName,
		},
	})
}

func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	var req dto.GetApplicationsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	if req.OpportunityID != nil {
		id, err := uuid.Parse(*req.OpportunityID)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
			return
		}

		applications, total, err := h.appService.GetApplicationsByOpportunity(c.Request.Context(), service.GetApplicationsByOpportunityDTO{
			Auth:          &auth,
			OpportunityID: id,
			Limit:         req.Limit,
			Offset:        req.Offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, service.ErrForbidden):
				c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			case errors.Is(err, service.ErrOpportunityNotFound):
				c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
			}
			return
		}

		resp := make([]*dto.ApplicationResponse, len(applications))
		for i, app := range applications {
			resp[i] = &dto.ApplicationResponse{
				ID:            app.ID,
				OpportunityID: app.OpportunityID,
				ApplicantID:   app.ApplicantID,
				Status:        app.Status,
				CreatedAt:     app.CreatedAt,
				UpdatedAt:     app.UpdatedAt,
				Opportunity: &dto.OpportunityInfo{
					ID:           app.OpportunityID,
					Title:        app.Opportunity.Title,
					Description:  app.Opportunity.Description,
					LocationCity: app.Opportunity.LocationCity,
				},
				Applicant: &dto.ApplicantInfo{
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
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid opportunity id"})
			return
		}

		applications, total, err := h.appService.GetApplicationsByApplicant(c.Request.Context(), service.GetApplicationsByApplicantDTO{
			Auth:        &auth,
			ApplicantID: id,
			Limit:       req.Limit,
			Offset:      req.Offset,
		})
		if err != nil {
			switch {
			case errors.Is(err, service.ErrForbidden):
				c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			case errors.Is(err, service.ErrApplicantNotFound):
				c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
			}
			return
		}

		resp := make([]*dto.ApplicationResponse, len(applications))
		for i, app := range applications {
			resp[i] = &dto.ApplicationResponse{
				ID:            app.ID,
				OpportunityID: app.OpportunityID,
				ApplicantID:   app.ApplicantID,
				Status:        app.Status,
				CreatedAt:     app.CreatedAt,
				UpdatedAt:     app.UpdatedAt,
				Opportunity: &dto.OpportunityInfo{
					ID:           app.OpportunityID,
					Title:        app.Opportunity.Title,
					Description:  app.Opportunity.Description,
					LocationCity: app.Opportunity.LocationCity,
				},
				Applicant: &dto.ApplicantInfo{
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

	c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "either opportunity_id or applicant_id is required"})
}
