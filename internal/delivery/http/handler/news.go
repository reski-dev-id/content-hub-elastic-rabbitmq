package handler

import (
	"net/http"

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

func (h *NewsHandler) Create(c *gin.Context) {
	var req entity.News

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.uc.Create(&req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "created")
}

func (h *NewsHandler) GetByID(c *gin.Context) {
	id := parseUint(c.Param("id"))

	data, err := h.uc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *NewsHandler) Update(c *gin.Context) {
	id := parseUint(c.Param("id"))

	var req entity.News

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req.ID = id

	if err := h.uc.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "updated")
}

func (h *NewsHandler) Delete(c *gin.Context) {
	id := parseUint(c.Param("id"))

	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "deleted")
}
