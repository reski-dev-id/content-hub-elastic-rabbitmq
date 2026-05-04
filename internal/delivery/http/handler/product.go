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

func (h *ProductHandler) GetAll(c *gin.Context) {
	data, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
