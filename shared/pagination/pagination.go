package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 10
	minPageSize     = 1
)

// Params represents pagination parameters extracted from a request.
type Params struct {
	Page     int
	PageSize int
}

// NewParams constructs Params while enforcing boundaries.
func NewParams(page, pageSize int) Params {
	if page < 1 {
		page = defaultPage
	}

	if pageSize < minPageSize {
		pageSize = defaultPageSize
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return Params{
		Page:     page,
		PageSize: pageSize,
	}
}

// FromQuery extracts pagination parameters from query string.
func FromQuery(ctx *gin.Context) Params {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", strconv.Itoa(defaultPage)))
	if err != nil {
		page = defaultPage
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", strconv.Itoa(defaultPageSize)))
	if err != nil {
		pageSize = defaultPageSize
	}

	return NewParams(page, pageSize)
}

// Offset returns the SQL offset value.
func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the SQL limit value.
func (p Params) Limit() int {
	return p.PageSize
}

// Metadata represents pagination information returned to clients.
type Metadata struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// NewMetadata builds Metadata from total rows and params.
func NewMetadata(totalItems int64, params Params) Metadata {
	totalPages := 0
	if params.PageSize > 0 && totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(params.PageSize)))
	}

	return Metadata{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
