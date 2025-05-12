package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"develapar-server/utils"
	"time"
)

type ArticleService interface {
	CreateArticle(payload model.Article) (model.Article, error)
	FindAll() ([]model.Article, error)
	UpdateArticle(id int, req dto.UpdateArticleRequest) (model.Article, error)
	FindBySlug(slug string) (model.Article, error)
	FindByUserId(userId int) ([]model.Article, error)
	FindByCategory(catId string) ([]model.Article, error)
	DeleteArticle(id int) error
}

type articleService struct {
	repo repository.ArticleRepository
}

// FindByCategory implements ArticleService.
func (a *articleService) FindByCategory(catId string) ([]model.Article, error) {
	return a.repo.GetArticleByCategory(catId)
}

// DeleteArticle implements ArticleService.
func (a *articleService) DeleteArticle(id int) error {
	return a.repo.DeleteArticle(id)
}

// FindByUserId implements ArticleService.
func (a *articleService) FindByUserId(userId int) ([]model.Article, error) {
	return a.repo.GetArticleByUserId(userId)
}

// FindBySlug implements ArticleService.
func (a *articleService) FindBySlug(slug string) (model.Article, error) {
	return a.repo.GetArticleBySlug(slug)
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
		slug := utils.GenerateSlug(*req.Slug)
		article.Slug = slug
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.CategoryID != nil {
		article.Category.Id = *req.CategoryID
	}

	return a.repo.UpdateArticle(article)
}

// CreateArticle implements ArticleService.
func (a *articleService) CreateArticle(payload model.Article) (model.Article, error) {
	slug := utils.GenerateSlug(payload.Slug)
	article := model.Article{

		Title:     payload.Title,
		Slug:      slug,
		Content:   payload.Content,
		User:      payload.User,
		Category:  payload.Category,
		Views:     payload.Views,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return a.repo.CreateArticle(article)
}

// FindAll implements ArticleService.
func (a *articleService) FindAll() ([]model.Article, error) {
	return a.repo.GetAll()
}

func NewArticleService(repository repository.ArticleRepository) ArticleService {
	return &articleService{repo: repository}
}
