package response

type SearchResponse struct {
	ID          interface{} `json:"id"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Description string      `json:"description,omitempty"`
	Content     string      `json:"content,omitempty"`
	Author      string      `json:"author,omitempty"`
	CategoryID  uint64      `json:"category_id"`
	Status      string      `json:"status"`
	Score       float64     `json:"score"`
}
