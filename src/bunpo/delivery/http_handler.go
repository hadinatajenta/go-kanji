package delivery

import (
	"github.com/gin-gonic/gin"

	"gobackend/shared/response"
	bunpointerfaces "gobackend/src/bunpo/interfaces"
)

// Handler exposes bunpo HTTP endpoints.
type Handler struct {
	service bunpointerfaces.Service
}

// NewHandler builds a Handler instance.
func NewHandler(service bunpointerfaces.Service) *Handler {
	return &Handler{service: service}
}

// Test responds with a simple success message.
func (h *Handler) Test(ctx *gin.Context) {
	message, err := h.service.Test(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "unable to execute bunpo test", err.Error())
		return
	}

	response.OK(ctx, message, gin.H{"message": message})
}
