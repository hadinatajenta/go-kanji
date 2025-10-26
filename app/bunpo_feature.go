package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	bunpodelivery "gobackend/src/bunpo/delivery"
	bunporoutes "gobackend/src/bunpo/routes"
	bunposervice "gobackend/src/bunpo/service"
)

// RegisterBunpoFeature wires the bunpo feature into the router.
func RegisterBunpoFeature(router gin.IRouter) error {
	if router == nil {
		return fmt.Errorf("register bunpo feature: router is nil")
	}

	service := bunposervice.NewService()
	handler := bunpodelivery.NewHandler(service)
	bunporoutes.Register(router, handler)

	return nil
}
