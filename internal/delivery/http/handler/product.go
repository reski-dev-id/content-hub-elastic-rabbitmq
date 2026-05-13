package handler

import (
	"net/http"

	"content-hub/internal/delivery/http/request"
	"content-hub/internal/delivery/http/response"
	"content-hub/internal/domain/entity"
	"content-hub/internal/domain/usecase"
	"content-hub/internal/helper"

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
			helper.FormatValidationError(err),
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

// SearchProducts godoc
// @Summary Search products
// @Description Search products with recommendations
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param category_id query int false "Category ID"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /v1/products/search [get]
func (h *ProductHandler) Search(c *gin.Context) {

	pagination := helper.ParsePagination(c)

	q := c.Query("q")

	var categoryID *uint64

	if cid := c.Query("category_id"); cid != "" {

		val := helper.ParseUint(cid)

		categoryID = &val
	}

	data, total, err := h.uc.Search(
		q,
		categoryID,
		pagination.Page,
		pagination.Limit,
	)

	if err != nil {

		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to search products",
			err.Error(),
		)

		return
	}

	response.SuccessWithMeta(
		c,
		http.StatusOK,
		"products search fetched successfully",
		data,
		helper.NewPagination(
			pagination.Page,
			pagination.Limit,
			total,
		),
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

	data, err := h.uc.GetByID(helper.ParseUint(id))

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
			helper.FormatValidationError(err),
		)

		return
	}

	product := entity.Product{
		ID:          helper.ParseUint(id),
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

	pagination := helper.ParsePagination(c)

	var categoryID *uint64

	if cid := c.Query("category_id"); cid != "" {

		val := helper.ParseUint(cid)

		categoryID = &val
	}

	data, total, err := h.uc.GetAll(
		pagination.Page,
		pagination.Limit,
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
		helper.NewPagination(
			pagination.Page,
			pagination.Limit,
			total,
		),
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

	id := helper.ParseUint(c.Param("id"))

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
