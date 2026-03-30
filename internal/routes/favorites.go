package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupFavoritesRoutes(
	protectedRouter *gin.RouterGroup,
	hnd *handler.FavoritesHandler,
) {
	protectedFavorites := protectedRouter.Group("/favorites")

	{
		protectedFavorites.GET("", hnd.List)
		protectedFavorites.POST("", hnd.Create)
		protectedFavorites.DELETE("/:id", hnd.Delete)
	}
}
