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
		slog.Warn("failed to extract auth context", "error", err, "path", c.Request.URL.Path)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	req := service.GetMeApplicantRequest{
		Auth: authCtx,
	}

	applicant, err := h.applicantService.GetMe(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found", "user_id", authCtx.UserID, "path", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		slog.Error("failed to get applicant", "error", err, "user_id", authCtx.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get applicant"})
		return
	}

	slog.Info("applicant retrieved successfully", "applicant_id", applicant.ID, "user_id", authCtx.UserID)
	c.JSON(http.StatusOK, toApplicantResponse(applicant, nil))
}

func (h *ApplicantHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("invalid applicant id", "id", idParam, "error", err)
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
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found", "applicant_id", applicantID, "path", c.Request.URL.Path)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied to applicant", "applicant_id", applicantID, "auth_user_id", func() string {
				if auth != nil {
					return auth.UserID.String()
				}
				return "none"
			}())
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		slog.Error("failed to get applicant by id", "error", err, "applicant_id", applicantID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get applicant"})
		return
	}

	slog.Info("applicant retrieved by id", "applicant_id", applicantID)
	c.JSON(http.StatusOK, toApplicantResponse(applicant, nil))
}

func (h *ApplicantHandler) List(c *gin.Context) {
	var req dto.ListApplicantsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("invalid query parameters for list applicants", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid query parameters"})
		return
	}

	serviceReq := service.ListApplicantsRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	applicants, err := h.applicantService.List(c.Request.Context(), serviceReq)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			slog.Warn("invalid input for list applicants", "error", err)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid limit or offset"})
			return
		}
		slog.Error("failed to list applicants", "error", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to list applicants"})
		return
	}

	applicantResponses := make([]dto.ApplicantInfo, len(applicants))
	for i, applicant := range applicants {
		applicantResponses[i] = dto.ApplicantInfo{
			ID:         applicant.ID,
			FirstName:  applicant.FirstName,
			SecondName: applicant.SecondName,
			LastName:   applicant.LastName,
		}
	}

	response := dto.ListApplicantsResponse{
		Applicants: applicantResponses,
		Total:      int64(len(applicantResponses)),
	}

	slog.Info("applicants listed successfully",
		"limit", req.Limit,
		"offset", req.Offset,
		"count", len(applicantResponses))

	c.JSON(http.StatusOK, response)
}

func (h *ApplicantHandler) Update(c *gin.Context) {
	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for update", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.UpdateApplicantRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		slog.Warn("invalid request body for update", "error", err, "user_id", authCtx.UserID)
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
		WorkExpirience: reqBody.WorkExpirience,
	}

	if err := h.applicantService.Update(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found for update", "user_id", authCtx.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied for update", "user_id", authCtx.UserID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			slog.Warn("invalid input for update", "error", err, "user_id", authCtx.UserID)
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		slog.Error("failed to update applicant", "error", err, "user_id", authCtx.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to update applicant"})
		return
	}

	slog.Info("applicant updated successfully", "user_id", authCtx.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "applicant updated successfully"})
}

func (h *ApplicantHandler) GetTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("invalid applicant id for get tags", "id", idParam, "error", err)
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
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found for get tags", "applicant_id", applicantID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied to get tags", "applicant_id", applicantID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		slog.Error("failed to get tags", "error", err, "applicant_id", applicantID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to get tags"})
		return
	}

	slog.Info("tags retrieved successfully", "applicant_id", applicantID, "tags_count", len(tags))
	c.JSON(http.StatusOK, tagsToResponse(tags))
}

func (h *ApplicantHandler) SetTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("invalid applicant id for set tags", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for set tags", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.TagsRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		slog.Warn("invalid request body for set tags", "error", err, "user_id", authCtx.UserID)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.SetApplicantTagsRequest{
		Auth:        authCtx,
		ApplicantID: applicantID,
		TagIDs:      reqBody.TagIDs,
	}

	if err := h.applicantService.SetTags(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found for set tags", "applicant_id", applicantID, "user_id", authCtx.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied for set tags", "applicant_id", applicantID, "user_id", authCtx.UserID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		slog.Error("failed to set tags", "error", err, "applicant_id", applicantID, "user_id", authCtx.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to set tags"})
		return
	}

	slog.Info("tags set successfully", "applicant_id", applicantID, "user_id", authCtx.UserID, "tags_count", len(reqBody.TagIDs))
	c.JSON(http.StatusOK, gin.H{"message": "tags set successfully"})
}

func (h *ApplicantHandler) AddTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("invalid applicant id for add tags", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for add tags", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.TagsRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		slog.Warn("invalid request body for add tags", "error", err, "user_id", authCtx.UserID)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.AddApplicantTagsRequest{
		Auth:        authCtx,
		ApplicantID: applicantID,
		TagIDs:      reqBody.TagIDs,
	}

	if err := h.applicantService.AddTags(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found for add tags", "applicant_id", applicantID, "user_id", authCtx.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied for add tags", "applicant_id", applicantID, "user_id", authCtx.UserID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		slog.Error("failed to add tags", "error", err, "applicant_id", applicantID, "user_id", authCtx.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to add tags"})
		return
	}

	slog.Info("tags added successfully", "applicant_id", applicantID, "user_id", authCtx.UserID, "tags_count", len(reqBody.TagIDs))
	c.JSON(http.StatusOK, gin.H{"message": "tags added successfully"})
}

func (h *ApplicantHandler) RemoveTags(c *gin.Context) {
	idParam := c.Param("id")
	applicantID, err := uuid.Parse(idParam)
	if err != nil {
		slog.Warn("invalid applicant id for remove tags", "id", idParam, "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid applicant id"})
		return
	}

	authCtx, err := extractAuthContext(c)
	if err != nil {
		slog.Warn("failed to extract auth context for remove tags", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var reqBody dto.TagsRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		slog.Warn("invalid request body for remove tags", "error", err, "user_id", authCtx.UserID)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	req := service.RemoveApplicantTagsRequest{
		Auth:        authCtx,
		ApplicantID: applicantID,
		TagIDs:      reqBody.TagIDs,
	}

	if err := h.applicantService.RemoveTags(c.Request.Context(), req); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			slog.Warn("applicant not found for remove tags", "applicant_id", applicantID, "user_id", authCtx.UserID)
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "applicant not found"})
			return
		}
		if errors.Is(err, service.ErrForbidden) {
			slog.Warn("access denied for remove tags", "applicant_id", applicantID, "user_id", authCtx.UserID)
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "access denied"})
			return
		}
		slog.Error("failed to remove tags", "error", err, "applicant_id", applicantID, "user_id", authCtx.UserID)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to remove tags"})
		return
	}

	slog.Info("tags removed successfully", "applicant_id", applicantID, "user_id", authCtx.UserID, "tags_count", len(reqBody.TagIDs))
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
		PrivacySetting: a.PrivacySetting,
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
