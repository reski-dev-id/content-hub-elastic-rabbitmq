package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page       int   `json:"page" example:"1"`
	Limit      int   `json:"limit" example:"10"`
	Total      int64 `json:"total" example:"100"`
	TotalPages int   `json:"total_pages" example:"10"`
}

type PaginationQuery struct {
	Page  int
	Limit int
}

func ParsePagination(c *gin.Context) PaginationQuery {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))

	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	return PaginationQuery{
		Page:  page,
		Limit: limit,
	}
}

func NewPagination(
	page int,
	limit int,
	total int64,
) Pagination {

	totalPages := int(total) / limit

	if int(total)%limit > 0 {
		totalPages++
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
