package request

type CreateProductRequest struct {
	CategoryID  uint64  `json:"category_id" binding:"required" example:"1"`
	Title       string  `json:"title" binding:"required,min=3" example:"MacBook Pro M4"`
	Slug        string  `json:"slug" binding:"required" example:"macbook-pro-m4"`
	Description string  `json:"description" example:"Laptop Apple terbaru"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"42000000"`
	Status      string  `json:"status" binding:"required,oneof=active inactive" example:"active"`
}

type UpdateProductRequest struct {
	CategoryID  uint64  `json:"category_id" binding:"required" example:"1"`
	Title       string  `json:"title" binding:"required,min=3" example:"MacBook Pro M4"`
	Slug        string  `json:"slug" binding:"required" example:"macbook-pro-m4"`
	Description string  `json:"description" example:"Laptop Apple terbaru"`
	Price       float64 `json:"price" binding:"required,gt=0" example:"42000000"`
	Status      string  `json:"status" binding:"required,oneof=active inactive" example:"active"`
}
