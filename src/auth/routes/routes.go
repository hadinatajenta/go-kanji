package routes

import (
    "github.com/gin-gonic/gin"

    "gobackend/src/auth/delivery"
)

// Register attaches the auth endpoints to the provided router.
func Register(router gin.IRoutes, handler *delivery.Handler) {
    router.GET("/auth/google/login", handler.GoogleLogin)
    router.GET("/auth/google/callback", handler.GoogleCallback)
    router.POST("/auth/logout", handler.Logout)
}

