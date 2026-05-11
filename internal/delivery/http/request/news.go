package request

type CreateNewsRequest struct {
	CategoryID uint64 `json:"category_id" binding:"required" example:"1"`
	Title      string `json:"title" binding:"required,min=3" example:"OpenAI Rilis GPT Baru"`
	Slug       string `json:"slug" binding:"required" example:"openai-rilis-gpt-baru"`
	Content    string `json:"content" binding:"required" example:"Model terbaru membawa peningkatan reasoning."`
	Author     string `json:"author" binding:"required" example:"Reski"`
	Status     string `json:"status" binding:"required,oneof=draft published" example:"published"`
}

type UpdateNewsRequest struct {
	CategoryID uint64 `json:"category_id" binding:"required" example:"1"`
	Title      string `json:"title" binding:"required,min=3" example:"OpenAI Rilis GPT Baru"`
	Slug       string `json:"slug" binding:"required" example:"openai-rilis-gpt-baru"`
	Content    string `json:"content" binding:"required" example:"Model terbaru membawa peningkatan reasoning."`
	Author     string `json:"author" binding:"required" example:"Reski"`
	Status     string `json:"status" binding:"required,oneof=draft published" example:"published"`
}
