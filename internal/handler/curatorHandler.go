package handler

import (
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req := service.GetMeRequest{
		Auth: auth,
	}

	curator, err := h.curatorService.GetMe(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	serviceReq := service.UpdateRequest{
		Auth:     auth,
		FullName: req.FullName,
	}

	updErr := h.curatorService.Update(c.Request.Context(), serviceReq)
	if updErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updErr.Error()})
		return
	}

	c.Status(http.StatusOK)
}
