package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupRecomendationRoutes(
	protectedRouter *gin.RouterGroup,
	hnd *handler.RecommendationHandler,
) {
	protectedRecomendations := protectedRouter.Group("/recomendations")
	{
		protectedRecomendations.POST("", hnd.Create)
		protectedRecomendations.GET("/my", hnd.GetMyRecomendations)
		protectedRecomendations.DELETE("/:id", hnd.Create)
	}
}
