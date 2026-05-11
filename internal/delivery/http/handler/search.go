package handler

import (
	"net/http"

	"content-hub/internal/delivery/http/response"
	"content-hub/internal/domain/usecase"
	"content-hub/internal/helper"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	uc usecase.SearchUsecase
}

func NewSearchHandler(
	uc usecase.SearchUsecase,
) *SearchHandler {
	return &SearchHandler{uc}
}

// SearchContent godoc
// @Summary Search products or news
// @Description Full-text search using Elasticsearch
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search keyword"
// @Param type query string true "Search type (product/news)"
// @Param category_id query int false "Category ID"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/search [get]
func (h *SearchHandler) Search(c *gin.Context) {

	q := c.Query("q")
	searchType := c.Query("type")

	pagination := helper.ParsePagination(c)

	var categoryID *uint64

	if cid := c.Query("category_id"); cid != "" {
		val := parseUint(cid)
		categoryID = &val
	}

	data, total, err := h.uc.Search(
		q,
		searchType,
		categoryID,
		pagination.Page,
		pagination.Limit,
	)

	if err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to search content",
			err.Error(),
		)

		return
	}

	response.SuccessWithMeta(
		c,
		http.StatusOK,
		"search completed successfully",
		data,
		helper.NewPagination(
			pagination.Page,
			pagination.Limit,
			total,
		),
	)
}
