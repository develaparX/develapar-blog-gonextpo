package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type ArticleService interface {
	CreateArticle(payload model.Article) (model.Article, error)
	FindAll() ([]model.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

// CreateArticle implements ArticleService.
func (a *articleService) CreateArticle(payload model.Article) (model.Article, error) {
	return a.repo.CreateArticle(payload)
}

// FindAll implements ArticleService.
func (a *articleService) FindAll() ([]model.Article, error) {
	return a.repo.GetAll()
}

func NewArticleService(repository repository.ArticleRepository) ArticleService {
	return &articleService{repo: repository}
}
