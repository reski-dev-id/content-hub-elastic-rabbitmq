package handler

import (
	"net/http"

	"content-hub/internal/domain/usecase"

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

func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	searchType := c.Query("type")

	page := parseIntDefault(
		c.Query("page"),
		1,
	)

	limit := parseIntDefault(
		c.Query("limit"),
		10,
	)

	var categoryID *uint64

	if cid := c.Query("category_id"); cid != "" {
		val := parseUint(cid)
		categoryID = &val
	}

	data, err := h.uc.Search(
		q,
		searchType,
		categoryID,
		page,
		limit,
	)

	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	c.JSON(http.StatusOK, data)
}
