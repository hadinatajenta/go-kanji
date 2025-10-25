package routes

import (
	"github.com/gin-gonic/gin"

	"gobackend/src/users/delivery"
)

func Register(router gin.IRoutes, handler *delivery.Handler) {
	router.GET("/api/users", handler.ListUsers)
}
