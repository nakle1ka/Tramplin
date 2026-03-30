package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupOpportunityRoutes(
	router *gin.RouterGroup,
	protectedRouter *gin.RouterGroup,
	hnd *handler.OpportunityHandler,
) {
	opportunities := router.Group("/opportunities")
	protectedOpportunities := protectedRouter.Group("/opportunities")

	{
		opportunities.GET("", hnd.List)
		opportunities.GET("/:id", hnd.GetByID)

		protectedOpportunities.POST("", hnd.Create)
		protectedOpportunities.PATCH("/:id", hnd.Update)
		protectedOpportunities.DELETE("/:id", hnd.Delete)

		protectedOpportunities.PATCH("/:id/moderation-status", hnd.UpdateModerationStatus)

		protectedOpportunities.POST("/:id/tags", hnd.AddTags)
		protectedOpportunities.DELETE("/:id/tags", hnd.RemoveTags)
	}
}
