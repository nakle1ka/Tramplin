package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupCuratorRoutes(
	protectedRouter *gin.RouterGroup,
	hnd *handler.CuratorHandler,
) {
	protectedCurators := protectedRouter.Group("/curators")
	{
		protectedCurators.GET("/me", hnd.GetMe)
		protectedCurators.PATCH("/me", hnd.Update)
	}
}
