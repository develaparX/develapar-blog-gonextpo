package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
)

type ArticleService interface {
	CreateArticle(payload model.Article) (model.Article, error)
	FindAll() ([]model.Article, error)
	UpdateArticle(id int, req dto.UpdateArticleRequest) (model.Article, error)
}

type articleService struct {
	repo repository.ArticleRepository
}

// UpdateArticle implements ArticleService.
func (a *articleService) UpdateArticle(id int, req dto.UpdateArticleRequest) (model.Article, error) {
	article, err := a.repo.GetArticleById(id)
	if err != nil {
		return model.Article{}, err

	}

	if req.Title != nil {
		article.Title = *req.Title
	}
	if req.Slug != nil {
		article.Slug = *req.Slug
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.CategoryID != nil {
		article.Category.Id = *req.CategoryID
	}

	// Step 3: Simpan ke DB
	return a.repo.UpdateArticle(article)
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
