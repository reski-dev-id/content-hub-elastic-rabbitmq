package handler

import (
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	uc usecase.ProductUsecase
}

func NewProductHandler(uc usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{uc}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req entity.Product

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

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	data, err := h.uc.GetByID(parseUint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entity.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req.ID = parseUint(id)

	if err := h.uc.Update(&req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "updated")
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 10)

	var categoryID *uint64
	if cid := c.Query("category_id"); cid != "" {
		val := parseUint(cid)
		categoryID = &val
	}

	data, err := h.uc.GetAll(page, limit, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := parseUint(c.Param("id"))

	if err := h.uc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "deleted")
}
