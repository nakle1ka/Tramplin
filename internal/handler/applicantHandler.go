package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/middleware"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type ApplicantHandler struct {
	service service.ApplicantService
}

func NewApplicantHandler(service service.ApplicantService) *ApplicantHandler {
	return &ApplicantHandler{
		service: service,
	}
}

func (h *ApplicantHandler) GetMe(c *gin.Context) {
	userIdStr, ok := c.Get(middleware.UserIdKey)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	userId := userIdStr.(uuid.UUID)

	applicant, err := h.service.GetMe(c.Request.Context(), userId)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.ApplicantResponse{
		ID:             applicant.ID,
		Email:          applicant.User.Email,
		UserID:         applicant.UserID,
		FirstName:      applicant.FirstName,
		SecondName:     applicant.SecondName,
		LastName:       applicant.LastName,
		University:     applicant.University,
		GraduationYear: applicant.GraduationYear,
		About:          applicant.About,
		PrivacySetting: int(applicant.PrivacySetting),
		CreatedAt:      applicant.CreatedAt,
		UpdatedAt:      applicant.UpdatedAt,
	}

	if len(applicant.Tags) > 0 {
		resp.Tags = make([]dto.TagResponse, len(applicant.Tags))
		for i, tag := range applicant.Tags {
			resp.Tags[i] = dto.TagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ApplicantHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	applicant, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := dto.ApplicantResponse{
		ID:             applicant.ID,
		Email:          applicant.User.Email,
		UserID:         applicant.UserID,
		FirstName:      applicant.FirstName,
		SecondName:     applicant.SecondName,
		LastName:       applicant.LastName,
		University:     applicant.University,
		GraduationYear: applicant.GraduationYear,
		About:          applicant.About,
		PrivacySetting: int(applicant.PrivacySetting),
		CreatedAt:      applicant.CreatedAt,
		UpdatedAt:      applicant.UpdatedAt,
	}

	if len(applicant.Tags) > 0 {
		resp.Tags = make([]dto.TagResponse, len(applicant.Tags))
		for i, tag := range applicant.Tags {
			resp.Tags[i] = dto.TagResponse{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ApplicantHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req dto.UpdateApplicantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	updateReq := service.UpdateApplicantRequest{
		FirstName:      req.FirstName,
		SecondName:     req.SecondName,
		LastName:       req.LastName,
		University:     req.University,
		GraduationYear: req.GraduationYear,
		About:          req.About,
		PrivacySetting: req.PrivacySetting,
	}

	if err := h.service.Update(c.Request.Context(), id, updateReq); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ApplicantHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *ApplicantHandler) AddTags(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	applicantID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req dto.TagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.AddTags(c.Request.Context(), applicantID, req.TagIDs); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tags added"})
}

func (h *ApplicantHandler) RemoveTags(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	applicantID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req dto.TagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.RemoveTags(c.Request.Context(), applicantID, req.TagIDs); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tags removed"})
}

func (h *ApplicantHandler) SetTags(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	applicantID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var req dto.TagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.SetTags(c.Request.Context(), applicantID, req.TagIDs); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		case errors.Is(err, service.ErrApplicantNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "applicant not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tags set"})
}
