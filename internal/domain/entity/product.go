package entity

type Product struct {
	ID          uint64
	CategoryID  uint64
	Title       string
	Slug        string
	Description string
	Price       float64
	Status      string
}
