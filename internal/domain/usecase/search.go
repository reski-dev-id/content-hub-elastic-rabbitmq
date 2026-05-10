package usecase

type SearchUsecase interface {
	Search(
		q string,
		searchType string,
		categoryID *uint64,
		page int,
		limit int,
	) (interface{}, error)
}
