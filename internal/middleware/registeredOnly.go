package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisteredOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, idExists := c.Get(UserIdKey)
		_, roleExists := c.Get(UserRoleKey)
		if !idExists || !roleExists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
