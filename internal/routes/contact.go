package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nakle1ka/Tramplin/internal/handler"
)

func SetupContactRoutes(
	protectedRouter *gin.RouterGroup,
	hnd *handler.ContactHandler,
) {
	protectedContacts := protectedRouter.Group("/contacts")
	{
		protectedContacts.POST("", hnd.Create)
		protectedContacts.GET("/friends", hnd.ListFriends)
		protectedContacts.GET("/sent", hnd.ListSentRequests)
		protectedContacts.GET("/received", hnd.ListReceivedRequests)
		protectedContacts.PUT("/:id/status", hnd.UpdateStatus)
		protectedContacts.DELETE("/:id", hnd.Delete)
	}
}
