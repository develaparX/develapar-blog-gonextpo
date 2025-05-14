package dto

type AssignTagsByNameDTO struct {
	ArticleID int      `json:"article_id" binding:"required"`
	Tags      []string `json:"tags" binding:"required"`
}
