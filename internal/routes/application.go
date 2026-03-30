package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupApplicationRoutes(
	protectedRouter *gin.RouterGroup,
	hnd *handler.ApplicationHandler,
) {
	protectedApplications := protectedRouter.Group("/applications")

	{
		protectedApplications.GET("/:id", hnd.GetApplicationByID)
		protectedApplications.GET("", hnd.GetApplications)
		protectedApplications.POST("", hnd.CreateApplication)
		protectedApplications.PUT("/:id/status", hnd.UpdateApplicationStatus)
		protectedApplications.DELETE("/:id", hnd.DeleteApplication)
	}
}
