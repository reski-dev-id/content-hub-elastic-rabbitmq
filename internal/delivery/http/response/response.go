package response

import "content-hub/internal/domain/entity"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ProductSearchResponse struct {
	Items           []entity.Product `json:"items"`
	Recommendations []entity.Product `json:"recommendations"`
}

type NewsSearchResponse struct {
	Items           []entity.News `json:"items"`
	Recommendations []entity.News `json:"recommendations"`
}
