package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/middleware"
	"github.com/nakle1ka/Tramplin/internal/model"
	"github.com/nakle1ka/Tramplin/internal/service"
)

func extractAuthContext(c *gin.Context) (service.AuthContext, error) {
	userIdVal, exists := c.Get(middleware.UserIdKey)
	if !exists {
		return service.AuthContext{}, errors.New("user id not found")
	}

	userId, ok := userIdVal.(uuid.UUID)
	if !ok {
		return service.AuthContext{}, errors.New("invalid user id type")
	}

	userRoleVal, exists := c.Get(middleware.UserRoleKey)
	if !exists {
		return service.AuthContext{}, errors.New("user role not found")
	}

	userRole, ok := userRoleVal.(model.Role)
	if !ok {
		return service.AuthContext{}, errors.New("invalid user role type")
	}

	return service.AuthContext{
		UserID: userId,
		Role:   userRole,
	}, nil
}
