package routes

import (
    "github.com/gin-gonic/gin"

    logdelivery "gobackend/src/logs/delivery"
)

// Register attaches log endpoints to the provided router.
func Register(router gin.IRoutes, handler *logdelivery.Handler) {
    router.GET("/api/users/logs", handler.ListLogs)
    router.GET("/api/users/:reference/logs", handler.ListLogsByUser)
}

