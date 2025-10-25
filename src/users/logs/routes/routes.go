package routes

import (
	"github.com/gin-gonic/gin"

	logdelivery "gobackend/src/users/logs/delivery"
)

// Register attaches log endpoints to the router.
func Register(router gin.IRoutes, handler *logdelivery.Handler) {
	router.GET("/api/users/logs", handler.ListLogs)
}
