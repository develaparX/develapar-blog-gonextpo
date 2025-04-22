package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type ArticleTagService interface {
	AssignTags(articleId int, tagId []int) error
	FindTagByArticleId(articleId int) ([]model.Tags, error)
	FindArticleByTagId(tagId int) ([]model.Article, error)
}

type articleTagService struct {
	repo repository.ArticleTagRepository
}

// AssignTags implements ArticleTagService.
func (a *articleTagService) AssignTags(articleId int, tagId []int) error {
	return a.repo.AssignTags(articleId, tagId)
}

// FindArticleByTagId implements ArticleTagService.
func (a *articleTagService) FindArticleByTagId(tagId int) ([]model.Article, error) {
	return a.repo.GetArticleByTagId(tagId)
}

// FindTagByArticleId implements ArticleTagService.
func (a *articleTagService) FindTagByArticleId(articleId int) ([]model.Tags, error) {
	return a.repo.GetTagsByArticleId(articleId)
}

func NewArticleTagService(repository repository.ArticleTagRepository) ArticleTagService {
	return &articleTagService{repo: repository}
}
