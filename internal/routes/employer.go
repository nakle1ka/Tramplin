package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupEmployerRoutes(router *gin.RouterGroup, protectedRouter *gin.RouterGroup, hnd *handler.EmployerHandler) {
	employers := router.Group("/employers")
	protectedEmployers := protectedRouter.Group("/employers")
	{
		employers.GET("/:id", hnd.GetByID)
		protectedEmployers.GET("", hnd.List)
		protectedEmployers.GET("/me", hnd.GetMe)
		protectedEmployers.PATCH("/me", hnd.Update)
		protectedEmployers.PUT("/:id/verification", hnd.UpdateVerificationStatus)
	}
}
