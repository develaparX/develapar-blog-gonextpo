package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"develapar-server/utils"
	"time"
)

type ArticleService interface {
	CreateArticleWithTags(req dto.CreateArticleRequest, userID int) (model.Article, error)
	FindAll() ([]model.Article, error)
	UpdateArticle(id int, req dto.UpdateArticleRequest) (model.Article, error)
	FindById(id int) (model.Article, error)
	FindBySlug(slug string) (model.Article, error)
	FindByUserId(userId int) ([]model.Article, error)
	FindByCategory(catId string) ([]model.Article, error)
	DeleteArticle(id int) error
}

type articleService struct {
	repo              repository.ArticleRepository
	articleTagService ArticleTagService
}

// FindById implements ArticleService.
func (a *articleService) FindById(id int) (model.Article, error) {
	return a.repo.GetArticleById(context.Background(), id)
}

// FindByCategory implements ArticleService.
func (a *articleService) FindByCategory(catId string) ([]model.Article, error) {
	return a.repo.GetArticleByCategory(context.Background(), catId)
}

// DeleteArticle implements ArticleService.
func (a *articleService) DeleteArticle(id int) error {
	return a.repo.DeleteArticle(context.Background(), id)
}

// FindByUserId implements ArticleService.
func (a *articleService) FindByUserId(userId int) ([]model.Article, error) {
	return a.repo.GetArticleByUserId(context.Background(), userId)
}

// FindBySlug implements ArticleService.
func (a *articleService) FindBySlug(slug string) (model.Article, error) {
	return a.repo.GetArticleBySlug(context.Background(), slug)
}

// UpdateArticle implements ArticleService.
func (a *articleService) UpdateArticle(id int, req dto.UpdateArticleRequest) (model.Article, error) {
	article, err := a.repo.GetArticleById(context.Background(), id)
	if err != nil {
		return model.Article{}, err
	}

	if req.Title != nil {
		article.Title = *req.Title
		// Generate new slug automatically when title is updated
		article.Slug = utils.GenerateSlug(*req.Title)
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.CategoryID != nil {
		article.Category.Id = *req.CategoryID
	}

	updatedArticle, err := a.repo.UpdateArticle(context.Background(), article)
	if err != nil {
		return model.Article{}, err
	}

	// Update tags if provided
	if len(req.Tags) > 0 {
		// Remove existing tags first, then assign new ones
		// This is a simple approach - in production you might want to be more selective
		err = a.assignTagsToArticle(updatedArticle.Id, req.Tags)
		if err != nil {
			// Log error but don't fail the update
			// In production, you might want to use database transactions
			return updatedArticle, err
		}
	}

	return updatedArticle, nil
}

// assignTagsToArticle is a helper method to assign tags to an article
// Uses ArticleTagService to avoid code duplication
func (a *articleService) assignTagsToArticle(articleId int, tagNames []string) error {
	return a.articleTagService.AsignTagsByName(articleId, tagNames)
}

// CreateArticleWithTags implements ArticleService.
func (a *articleService) CreateArticleWithTags(req dto.CreateArticleRequest, userID int) (model.Article, error) {
	// Generate slug automatically from title
	slug := utils.GenerateSlug(req.Title)
	
	// Create article first
	article := model.Article{
		Title:     req.Title,
		Slug:      slug,
		Content:   req.Content,
		User:      model.User{Id: userID},
		Category:  model.Category{Id: req.CategoryID},
		Views:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdArticle, err := a.repo.CreateArticle(context.Background(), article)
	if err != nil {
		return model.Article{}, err
	}

	// Assign tags if provided
	if len(req.Tags) > 0 {
		err = a.assignTagsToArticle(createdArticle.Id, req.Tags)
		if err != nil {
			// If tag assignment fails, we could either:
			// 1. Delete the created article (rollback)
			// 2. Return the article without tags
			// For now, we'll return the article without tags and log the error
			// In production, you might want to use database transactions
			return createdArticle, err
		}
	}

	return createdArticle, nil
}



// FindAll implements ArticleService.
func (a *articleService) FindAll() ([]model.Article, error) {
	return a.repo.GetAll(context.Background())
}

func NewArticleService(repository repository.ArticleRepository, articleTagService ArticleTagService) ArticleService {
	return &articleService{
		repo:              repository,
		articleTagService: articleTagService,
	}
}
