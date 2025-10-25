package delivery

import (
	"github.com/gin-gonic/gin"

	"gobackend/shared/response"
	userinterfaces "gobackend/src/users/interfaces"
)

// Handler exposes HTTP handlers for the user feature.
type Handler struct {
	service userinterfaces.UserService
}

// NewHandler builds a new Handler.
func NewHandler(service userinterfaces.UserService) *Handler {
	return &Handler{service: service}
}

// ListUsers returns all registered users.
func (h *Handler) ListUsers(ctx *gin.Context) {
	users, err := h.service.ListUsers(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "failed to list users", err.Error())
		return
	}

	response.OK(ctx, "users retrieved successfully", gin.H{
		"users": users,
		"count": len(users),
	})
}
