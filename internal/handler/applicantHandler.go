package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type ApplicantHandler struct {
	applicantService service.ApplicantService
}

func NewApplicantHandler(applicantService service.ApplicantService) *ApplicantHandler {
	return &ApplicantHandler{
		applicantService: applicantService,
	}
}

func (h *ApplicantHandler) GetMe(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	req := service.GetMeApplicantRequest{
		Auth: authCtx,
	}

	applicant, err := h.applicantService.GetMe(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get applicant"})
		return
	}

	c.JSON(http.StatusOK, toApplicantResponse(applicant, nil))
}

func (h *ApplicantHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	var auth *service.AuthContext
	authCtx, err := extractAuthContext(c)
	if err == nil {
		auth = &authCtx
	}

	req := service.GetApplicantByIDRequest{
		Auth: auth,
		ID:   applicantID,
	}

	applicant, err := h.applicantService.GetByID(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get applicant"})
		return
	}

	c.JSON(http.StatusOK, toApplicantResponse(applicant, nil))
}

func (h *ApplicantHandler) Update(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.UpdateApplicantRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.UpdateApplicantRequest{
		Auth:           authCtx,
		FirstName:      reqBody.FirstName,
		SecondName:     reqBody.SecondName,
		LastName:       reqBody.LastName,
		University:     reqBody.University,
		GraduationYear: reqBody.GraduationYear,
		About:          reqBody.About,
		PrivacySetting: reqBody.PrivacySetting,
	}

	if err := h.applicantService.Update(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update applicant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "applicant updated successfully"})
}

func (h *ApplicantHandler) GetTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	var auth *service.AuthContext
	authCtx, err := extractAuthContext(c)
	if err == nil {
		auth = &authCtx
	}

	req := service.GetApplicantTagsRequest{
		Auth:        auth,
		ApplicantID: applicantID,
	}

	tags, err := h.applicantService.GetTags(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get tags"})
		return
	}

	c.JSON(http.StatusOK, tagsToResponse(tags))
}

func (h *ApplicantHandler) SetTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.TagsRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.SetApplicantTagsRequest{
		Auth:        authCtx,
		ApplicantID: applicantID,
		TagIDs:      reqBody.TagIDs,
	}

	if err := h.applicantService.SetTags(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to set tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tags set successfully"})
}

func (h *ApplicantHandler) AddTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.TagsRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.AddApplicantTagsRequest{
		Auth:        authCtx,
		ApplicantID: applicantID,
		TagIDs:      reqBody.TagIDs,
	}

	if err := h.applicantService.AddTags(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to add tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tags added successfully"})
}

func (h *ApplicantHandler) RemoveTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.TagsRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.RemoveApplicantTagsRequest{
		Auth:        authCtx,
		ApplicantID: applicantID,
		TagIDs:      reqBody.TagIDs,
	}

	if err := h.applicantService.RemoveTags(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrApplicantNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to remove tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tags removed successfully"})
}

func toApplicantResponse(a *model.Applicant, tags []*model.Tag) dto.ApplicantResponse {
	resp := dto.ApplicantResponse{
		ID:             a.ID,
		Email:          a.User.Email,
		UserID:         a.UserID,
		FirstName:      a.FirstName,
		SecondName:     a.SecondName,
		LastName:       a.LastName,
		University:     a.University,
		GraduationYear: a.GraduationYear,
		About:          a.About,
		PrivacySetting: int(a.PrivacySetting),
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}

	if tags != nil {
		resp.Tags = tagsToResponse(tags)
	}

	return resp
}

func tagsToResponse(tags []*model.Tag) []dto.TagResponse {
	result := make([]dto.TagResponse, len(tags))
	for i, tag := range tags {
		result[i] = dto.TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}
	return result
}
