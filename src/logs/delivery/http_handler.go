package delivery

import (
	"github.com/gin-gonic/gin"

	"gobackend/shared/identity"
	"gobackend/shared/pagination"
	"gobackend/shared/response"
	loginterfaces "gobackend/src/logs/interfaces"
)

// Handler exposes endpoints for user logs.
type Handler struct {
	service    loginterfaces.Service
	refEncoder *identity.UserReferenceEncoder
}

// NewHandler constructs a Handler.
func NewHandler(service loginterfaces.Service, refEncoder *identity.UserReferenceEncoder) *Handler {
	return &Handler{service: service, refEncoder: refEncoder}
}

// ListLogs returns paginated user log entries.
func (h *Handler) ListLogs(ctx *gin.Context) {
	params := pagination.FromQuery(ctx)

	reference := ctx.Query("reference")
	var userID *int64
	if reference != "" {
		decoded, err := h.refEncoder.Decode(reference)
		if err != nil {
			response.BadRequest(ctx, "invalid user reference", err.Error())
			return
		}
		userID = &decoded
	}

	logs, total, err := h.service.ListLogs(ctx.Request.Context(), params, userID)
	if err != nil {
		response.InternalError(ctx, "failed to fetch user logs", err.Error())
		return
	}

	meta := pagination.NewMetadata(total, params)
	response.Paginated(ctx, "user logs retrieved successfully", logs, meta)
}

// ListLogsByUser returns paginated log entries for a specific user reference.
func (h *Handler) ListLogsByUser(ctx *gin.Context) {
	params := pagination.FromQuery(ctx)

	reference := ctx.Param("reference")
	decoded, err := h.refEncoder.Decode(reference)
	if err != nil {
		response.BadRequest(ctx, "invalid user reference", err.Error())
		return
	}

	logs, total, err := h.service.ListLogs(ctx.Request.Context(), params, &decoded)
	if err != nil {
		response.InternalError(ctx, "failed to fetch user logs", err.Error())
		return
	}

	meta := pagination.NewMetadata(total, params)
	response.Paginated(ctx, "user logs retrieved successfully", logs, meta)
}
