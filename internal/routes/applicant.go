package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupApplicantRoutes(
	router *gin.RouterGroup,
	protectedRouter *gin.RouterGroup,
	hnd *handler.ApplicantHandler,
) {
	applicant := router.Group("/applicants")
	protectedApplicant := protectedRouter.Group("/applicants")

	{
		applicant.GET("", hnd.List)
		applicant.GET("/:id", hnd.GetByID)
		protectedApplicant.GET("/me", hnd.GetMe)
		protectedApplicant.PATCH("/me", hnd.Update)
		protectedApplicant.POST("/:id/tags", hnd.AddTags)
		protectedApplicant.DELETE("/:id/tags", hnd.RemoveTags)
	}
}
