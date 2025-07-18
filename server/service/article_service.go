package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"develapar-server/utils"
	"fmt"
	"time"
)

type ArticleService interface {
	CreateArticleWithTags(ctx context.Context, req dto.CreateArticleRequest, userID int) (model.Article, error)
	FindAll(ctx context.Context) ([]model.Article, error)
	FindAllWithPagination(ctx context.Context, page, limit int) (PaginationResult, error)
	UpdateArticle(ctx context.Context, id int, req dto.UpdateArticleRequest) (model.Article, error)
	FindById(ctx context.Context, id int) (model.Article, error)
	FindBySlug(ctx context.Context, slug string) (model.Article, error)
	FindByUserId(ctx context.Context, userId int) ([]model.Article, error)
	FindByUserIdWithPagination(ctx context.Context, userId, page, limit int) (PaginationResult, error)
	FindByCategory(ctx context.Context, catId string) ([]model.Article, error)
	FindByCategoryWithPagination(ctx context.Context, catId string, page, limit int) (PaginationResult, error)
	DeleteArticle(ctx context.Context, id int) error
}

type articleService struct {
	repo              repository.ArticleRepository
	articleTagService ArticleTagService
	paginationService PaginationService
	validationService ValidationService
}

// FindById implements ArticleService.
func (a *articleService) FindById(ctx context.Context, id int) (model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Article{}, ctx.Err()
	default:
	}

	// Validate ID
	if id <= 0 {
		return model.Article{}, fmt.Errorf("article ID must be greater than 0")
	}

	// Get article from repository with context
	article, err := a.repo.GetArticleById(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, fmt.Errorf("failed to fetch article: %v", err)
	}

	return article, nil
}

// FindByCategory implements ArticleService.
func (a *articleService) FindByCategory(ctx context.Context, catId string) ([]model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate category ID
	if catId == "" {
		return nil, fmt.Errorf("category ID is required")
	}

	// Get articles by category from repository with context
	articles, err := a.repo.GetArticleByCategory(ctx, catId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch articles by category: %v", err)
	}

	return articles, nil
}

// DeleteArticle implements ArticleService.
func (a *articleService) DeleteArticle(ctx context.Context, id int) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate ID
	if id <= 0 {
		return fmt.Errorf("article ID must be greater than 0")
	}

	// Delete article from repository with context
	err := a.repo.DeleteArticle(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to delete article: %v", err)
	}

	return nil
}

// FindByUserId implements ArticleService.
func (a *articleService) FindByUserId(ctx context.Context, userId int) ([]model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate user ID
	if userId <= 0 {
		return nil, fmt.Errorf("user ID must be greater than 0")
	}

	// Get articles by user from repository with context
	articles, err := a.repo.GetArticleByUserId(ctx, userId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch articles by user: %v", err)
	}

	return articles, nil
}

// FindBySlug implements ArticleService.
func (a *articleService) FindBySlug(ctx context.Context, slug string) (model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Article{}, ctx.Err()
	default:
	}

	// Validate slug
	if slug == "" {
		return model.Article{}, fmt.Errorf("slug is required")
	}

	// Get article by slug from repository with context
	article, err := a.repo.GetArticleBySlug(ctx, slug)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, fmt.Errorf("failed to fetch article by slug: %v", err)
	}

	return article, nil
}

// UpdateArticle implements ArticleService.
func (a *articleService) UpdateArticle(ctx context.Context, id int, req dto.UpdateArticleRequest) (model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Article{}, ctx.Err()
	default:
	}

	// Validate ID
	if id <= 0 {
		return model.Article{}, fmt.Errorf("article ID must be greater than 0")
	}

	// Get existing article with context
	article, err := a.repo.GetArticleById(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, fmt.Errorf("failed to fetch article for update: %v", err)
	}

	// Update fields if provided
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

	// Validate updated article data
	if validationErr := a.validationService.ValidateArticle(ctx, article); validationErr != nil {
		return model.Article{}, validationErr
	}

	// Check context cancellation before update
	select {
	case <-ctx.Done():
		return model.Article{}, ctx.Err()
	default:
	}

	// Update article in repository with context
	updatedArticle, err := a.repo.UpdateArticle(ctx, article)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, fmt.Errorf("failed to update article: %v", err)
	}

	// Update tags if provided
	if len(req.Tags) > 0 {
		// Remove existing tags first, then assign new ones
		// This is a simple approach - in production you might want to be more selective
		err = a.assignTagsToArticle(ctx, updatedArticle.Id, req.Tags)
		if err != nil {
			// Log error but don't fail the update
			// In production, you might want to use database transactions
			return updatedArticle, fmt.Errorf("failed to update article tags: %v", err)
		}
	}

	return updatedArticle, nil
}

// assignTagsToArticle is a helper method to assign tags to an article
// Uses ArticleTagService to avoid code duplication
func (a *articleService) assignTagsToArticle(ctx context.Context, articleId int, tagNames []string) error {
	return a.articleTagService.AsignTagsByName(ctx, articleId, tagNames)
}

// CreateArticleWithTags implements ArticleService.
func (a *articleService) CreateArticleWithTags(ctx context.Context, req dto.CreateArticleRequest, userID int) (model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Article{}, ctx.Err()
	default:
	}

	// Validate user ID
	if userID <= 0 {
		return model.Article{}, fmt.Errorf("valid user ID is required")
	}

	// Generate slug automatically from title
	slug := utils.GenerateSlug(req.Title)
	
	// Create article object
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

	// Validate article data using validation service
	if validationErr := a.validationService.ValidateArticle(ctx, article); validationErr != nil {
		return model.Article{}, validationErr
	}

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.Article{}, ctx.Err()
	default:
	}

	// Create article in repository with context
	createdArticle, err := a.repo.CreateArticle(ctx, article)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, fmt.Errorf("failed to create article: %v", err)
	}

	// Assign tags if provided
	if len(req.Tags) > 0 {
		err = a.assignTagsToArticle(ctx, createdArticle.Id, req.Tags)
		if err != nil {
			// If tag assignment fails, we could either:
			// 1. Delete the created article (rollback)
			// 2. Return the article without tags
			// For now, we'll return the article without tags and log the error
			// In production, you might want to use database transactions
			return createdArticle, fmt.Errorf("failed to assign tags to article: %v", err)
		}
	}

	return createdArticle, nil
}



// FindAll implements ArticleService.
func (a *articleService) FindAll(ctx context.Context) ([]model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Get all articles from repository with context
	articles, err := a.repo.GetAll(ctx)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch articles: %v", err)
	}

	return articles, nil
}

// FindAllWithPagination implements ArticleService with pagination support
func (a *articleService) FindAllWithPagination(ctx context.Context, page, limit int) (PaginationResult, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := a.paginationService.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, fmt.Errorf("pagination validation failed: %v", err)
	}

	// Get paginated articles from repository
	articles, total, repoErr := a.repo.GetAllWithPagination(ctx, query.Offset, query.Limit)
	if repoErr != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, fmt.Errorf("failed to fetch articles: %v", repoErr)
	}

	// Create pagination result
	result, paginationErr := a.paginationService.Paginate(ctx, articles, total, query)
	if paginationErr != nil {
		return PaginationResult{}, fmt.Errorf("failed to create pagination result: %v", paginationErr)
	}

	return result, nil
}

// FindByUserIdWithPagination implements ArticleService with pagination support for user articles
func (a *articleService) FindByUserIdWithPagination(ctx context.Context, userId, page, limit int) (PaginationResult, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := a.paginationService.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, fmt.Errorf("pagination validation failed: %v", err)
	}

	// Get paginated articles by user from repository
	articles, total, repoErr := a.repo.GetArticleByUserIdWithPagination(ctx, userId, query.Offset, query.Limit)
	if repoErr != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, fmt.Errorf("failed to fetch user articles: %v", repoErr)
	}

	// Create pagination result
	result, paginationErr := a.paginationService.Paginate(ctx, articles, total, query)
	if paginationErr != nil {
		return PaginationResult{}, fmt.Errorf("failed to create pagination result: %v", paginationErr)
	}

	return result, nil
}

// FindByCategoryWithPagination implements ArticleService with pagination support for category articles
func (a *articleService) FindByCategoryWithPagination(ctx context.Context, catId string, page, limit int) (PaginationResult, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := a.paginationService.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, fmt.Errorf("pagination validation failed: %v", err)
	}

	// Get paginated articles by category from repository
	articles, total, repoErr := a.repo.GetArticleByCategoryWithPagination(ctx, catId, query.Offset, query.Limit)
	if repoErr != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, fmt.Errorf("failed to fetch category articles: %v", repoErr)
	}

	// Create pagination result
	result, paginationErr := a.paginationService.Paginate(ctx, articles, total, query)
	if paginationErr != nil {
		return PaginationResult{}, fmt.Errorf("failed to create pagination result: %v", paginationErr)
	}

	return result, nil
}

func NewArticleService(repository repository.ArticleRepository, articleTagService ArticleTagService, paginationService PaginationService, validationService ValidationService) ArticleService {
	return &articleService{
		repo:              repository,
		articleTagService: articleTagService,
		paginationService: paginationService,
		validationService: validationService,
	}
}
