package handler

import (
	"net/http"

	"content-hub/internal/delivery/http/response"
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	uc usecase.NewsUsecase
}

func NewNewsHandler(
	uc usecase.NewsUsecase,
) *NewsHandler {
	return &NewsHandler{uc}
}

// CreateNews godoc
// @Summary Create news
// @Description Create new news
// @Tags news
// @Accept json
// @Produce json
// @Param request body entity.News true "News payload"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/news [post]
func (h *NewsHandler) Create(c *gin.Context) {

	var req entity.News

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)

		return
	}

	if err := h.uc.Create(&req); err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to create news",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"news created successfully",
		req,
	)
}

// GetNewsByID godoc
// @Summary Get news detail
// @Description Get news by ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/news/{id} [get]
func (h *NewsHandler) GetByID(c *gin.Context) {

	id := parseUint(c.Param("id"))

	data, err := h.uc.GetByID(id)

	if err != nil {

		response.Error(
			c,
			http.StatusNotFound,
			"news not found",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"news fetched successfully",
		data,
	)
}

// GetNews godoc
// @Summary Get news list
// @Description Get news with pagination
// @Tags news
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param category_id query int false "Category ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/news [get]
func (h *NewsHandler) GetAll(c *gin.Context) {

	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 10)

	var categoryID *uint64

	if cid := c.Query("category_id"); cid != "" {
		val := parseUint(cid)
		categoryID = &val
	}

	data, err := h.uc.GetAll(
		page,
		limit,
		categoryID,
	)

	if err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to fetch news",
			err.Error(),
		)

		return
	}

	response.SuccessWithMeta(
		c,
		http.StatusOK,
		"news fetched successfully",
		data,
		gin.H{
			"page":  page,
			"limit": limit,
		},
	)
}

// UpdateNews godoc
// @Summary Update news
// @Description Update news by ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Param request body entity.News true "News payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/news/{id} [put]
func (h *NewsHandler) Update(c *gin.Context) {

	id := parseUint(c.Param("id"))

	var req entity.News

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
			err.Error(),
		)

		return
	}

	req.ID = id

	if err := h.uc.Update(&req); err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to update news",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"news updated successfully",
		req,
	)
}

// DeleteNews godoc
// @Summary Delete news
// @Description Delete news by ID
// @Tags news
// @Accept json
// @Produce json
// @Param id path int true "News ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/news/{id} [delete]
func (h *NewsHandler) Delete(c *gin.Context) {

	id := parseUint(c.Param("id"))

	if err := h.uc.Delete(id); err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to delete news",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"news deleted successfully",
		nil,
	)
}
