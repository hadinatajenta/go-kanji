package delivery

import (
	"github.com/gin-gonic/gin"

	"gobackend/shared/pagination"
	"gobackend/shared/response"
	loginterfaces "gobackend/src/logs/interfaces"
)

// Handler exposes endpoints for user logs.
type Handler struct {
	service loginterfaces.Service
}

// NewHandler constructs a Handler.
func NewHandler(service loginterfaces.Service) *Handler {
	return &Handler{service: service}
}

// ListLogs returns paginated user log entries.
func (h *Handler) ListLogs(ctx *gin.Context) {
	params := pagination.FromQuery(ctx)

	logs, total, err := h.service.ListLogs(ctx.Request.Context(), params)
	if err != nil {
		response.InternalError(ctx, "failed to fetch user logs", err.Error())
		return
	}

	meta := pagination.NewMetadata(total, params)
	response.Paginated(ctx, "user logs retrieved successfully", logs, meta)
}
