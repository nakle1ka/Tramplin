package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/pkg/auth"
)

const (
	UserIdKey   = "x-user-id"
	UserRoleKey = "x-user-role"
)

func ParseJWT(tm auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Next()
			return
		}
		accessToken := strings.TrimPrefix(header, "Bearer ")

		claims, err := tm.ValidateToken(accessToken)
		if err != nil {
			c.Next()
			return
		}

		userId, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserIdKey, userId)
		c.Set(UserRoleKey, claims.UserRole)
		c.Next()
	}
}
