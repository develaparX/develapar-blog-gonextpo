package dto

type UpdateArticleRequest struct {
	Title      *string `json:"title"`
	Slug       *string `json:"slug"`
	Content    *string `json:"content"`
	CategoryID *int    `json:"category_id"`
}

type UpdateCategoryRequest struct {
	Name *string `json:"name"`
}
