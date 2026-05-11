package handler

import (
	"net/http"

	"content-hub/internal/delivery/http/request"
	"content-hub/internal/delivery/http/response"
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	uc usecase.ProductUsecase
}

func NewProductHandler(
	uc usecase.ProductUsecase,
) *ProductHandler {
	return &ProductHandler{uc}
}

// CreateProduct godoc
// @Summary Create product
// @Description Create new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body request.CreateProductRequest true "Product payload"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/products [post]
func (h *ProductHandler) Create(c *gin.Context) {

	var req request.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			"validation failed",
			err.Error(),
		)

		return
	}

	product := entity.Product{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		Price:       req.Price,
		Status:      req.Status,
	}

	if err := h.uc.Create(&product); err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to create product",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"product created successfully",
		product,
	)
}

// GetProductByID godoc
// @Summary Get product detail
// @Description Get product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {

	id := c.Param("id")

	data, err := h.uc.GetByID(parseUint(id))

	if err != nil {

		response.Error(
			c,
			http.StatusNotFound,
			"product not found",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"product fetched successfully",
		data,
	)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body request.UpdateProductRequest true "Product payload"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {

	id := c.Param("id")

	var req request.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			"validation failed",
			err.Error(),
		)

		return
	}

	product := entity.Product{
		ID:          parseUint(id),
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		Price:       req.Price,
		Status:      req.Status,
	}

	if err := h.uc.Update(&product); err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to update product",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"product updated successfully",
		product,
	)
}

// GetProducts godoc
// @Summary Get products
// @Description Get product list with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param category_id query int false "Category ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {

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
			"failed to fetch products",
			err.Error(),
		)

		return
	}

	response.SuccessWithMeta(
		c,
		http.StatusOK,
		"products fetched successfully",
		data,
		gin.H{
			"page":  page,
			"limit": limit,
		},
	)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {

	id := parseUint(c.Param("id"))

	if err := h.uc.Delete(id); err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to delete product",
			err.Error(),
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"product deleted successfully",
		nil,
	)
}
