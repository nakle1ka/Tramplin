package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
	"github.com/nakle1ka/Tramplin/internal/middleware"
	"github.com/nakle1ka/Tramplin/internal/pkg/auth"
)

func SetupApplicantRoutes(
	router *gin.RouterGroup,
	tm auth.TokenManager,
	hnd *handler.ApplicantHandler,
) {
	applicant := router.Group("/applicants")

	protected := applicant.Group("/")
	protected.Use(middleware.JWTAuth(tm))

	{
		applicant.GET("/:id", hnd.GetByID)
		protected.GET("/me", hnd.GetMe)
		protected.PATCH("/:id", hnd.Update)
		protected.DELETE("/:id", hnd.Delete)
		protected.POST("/:id/tags", hnd.AddTags)
		protected.DELETE("/:id/tags", hnd.RemoveTags)
	}
}
