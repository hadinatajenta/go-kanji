package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gobackend/shared/pagination"
)

func statusLabel(code int) string {
	if code >= 200 && code < 300 {
		return "success"
	}
	return "error"
}

// JSON writes a standardised API response envelope.
func JSON(ctx *gin.Context, status int, message string, data interface{}, errs interface{}) {
	if status == http.StatusNoContent {
		ctx.Status(http.StatusNoContent)
		return
	}

	payload := Envelope{
		Status:  statusLabel(status),
		Message: message,
	}

	if data != nil {
		payload.Data = data
	}

	if errs != nil {
		payload.Errors = errs
	}

	ctx.JSON(status, payload)
}

// OK returns a 200 response.
func OK(ctx *gin.Context, message string, data interface{}) {
	JSON(ctx, http.StatusOK, message, data, nil)
}

// Paginated returns a 200 response with pagination metadata.
func Paginated(ctx *gin.Context, message string, items interface{}, meta pagination.Metadata) {
	payload := gin.H{
		"items": items,
		"meta":  meta,
	}
	JSON(ctx, http.StatusOK, message, payload, nil)
}

// Created returns a 201 response.
func Created(ctx *gin.Context, message string, data interface{}) {
	JSON(ctx, http.StatusCreated, message, data, nil)
}

// NoContent returns a 204 response with no body.
func NoContent(ctx *gin.Context) {
	JSON(ctx, http.StatusNoContent, "", nil, nil)
}

// BadRequest returns a 400 response.
func BadRequest(ctx *gin.Context, message string, errs interface{}) {
	JSON(ctx, http.StatusBadRequest, message, nil, errs)
}

// Unauthorized returns a 401 response.
func Unauthorized(ctx *gin.Context, message string) {
	JSON(ctx, http.StatusUnauthorized, message, nil, nil)
}

// Forbidden returns a 403 response.
func Forbidden(ctx *gin.Context, message string) {
	JSON(ctx, http.StatusForbidden, message, nil, nil)
}

// NotFound returns a 404 response.
func NotFound(ctx *gin.Context, message string) {
	JSON(ctx, http.StatusNotFound, message, nil, nil)
}

// InternalError returns a 500 response.
func InternalError(ctx *gin.Context, message string, errs interface{}) {
	JSON(ctx, http.StatusInternalServerError, message, nil, errs)
}
