package routes

import (
	"github.com/gin-gonic/gin"

	bunpodelivery "gobackend/src/bunpo/delivery"
)

// Register mounts bunpo endpoints on the router.
func Register(router gin.IRoutes, handler *bunpodelivery.Handler) {
	router.GET("/bunpo/test", handler.Test)
}
