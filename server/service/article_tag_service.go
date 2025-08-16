package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ArticleTagService interface {
	AssignTags(ctx context.Context, articleId uuid.UUID, tagId []uuid.UUID) error
	AsignTagsByName(ctx context.Context, articleId uuid.UUID, tagNames []string) error
	FindTagByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Tags, error)
	FindArticleByTagId(ctx context.Context, tagId uuid.UUID) ([]model.Article, error)
	RemoveTagFromArticle(ctx context.Context, articleId, tagId uuid.UUID) error
}

type articleTagService struct {
	articleTagRepo    repository.ArticleTagRepository
	tagRepo           repository.TagRepository
	validationService ValidationService
}

// RemoveTagFromArticle implements ArticleTagService.
func (a *articleTagService) RemoveTagFromArticle(ctx context.Context, articleId, tagId uuid.UUID) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate IDs
	if articleId == uuid.Nil {
		return fmt.Errorf("article ID must be greater than 0")
	}
	if tagId == uuid.Nil {
		return fmt.Errorf("tag ID must be greater than 0")
	}

	// Remove tag from article with context
	err := a.articleTagRepo.RemoveTagFromArticle(ctx, articleId, tagId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to remove tag from article: %v", err)
	}

	return nil
}

// AsignTagsByName implements ArticleTagService.
func (a *articleTagService) AsignTagsByName(ctx context.Context, articleId uuid.UUID, tagNames []string) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate inputs
	if articleId == uuid.Nil {
		return fmt.Errorf("article ID must be greater than 0")
	}
	if len(tagNames) == 0 {
		return fmt.Errorf("tag names are required")
	}

	var tagIds []uuid.UUID

	for _, tagName := range tagNames {
		// Check context cancellation in loop
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Validate and normalize tag name
		tagName = strings.ToLower(strings.TrimSpace(tagName))
		if tagName == "" {
			continue // Skip empty tag names
		}

		// Try to get existing tag with context
		tag, err := a.tagRepo.GetTagByName(ctx, tagName)
		if err != nil {
			// Check if context was cancelled during repository operation
			if ctx.Err() != nil {
				return ctx.Err()
			}

			// Tag doesn't exist, create new one with context
			newTag, createErr := a.tagRepo.CreateTag(ctx, model.Tags{Name: tagName})
			if createErr != nil {
				// Check if context was cancelled during repository operation
				if ctx.Err() != nil {
					return ctx.Err()
				}
				return fmt.Errorf("failed to create tag '%s': %v", tagName, createErr)
			}
			tagIds = append(tagIds, newTag.Id)
		} else {
			tagIds = append(tagIds, tag.Id)
		}
	}

	// Check context cancellation before final assignment
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Assign tags to article with context
	err := a.articleTagRepo.AssignTags(ctx, articleId, tagIds)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to assign tags to article: %v", err)
	}

	return nil
}

// AssignTags implements ArticleTagService.
func (a *articleTagService) AssignTags(ctx context.Context, articleId uuid.UUID, tagId []uuid.UUID) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate inputs
	if articleId == uuid.Nil {
		return fmt.Errorf("article ID must be greater than 0")
	}
	if len(tagId) == 0 {
		return fmt.Errorf("tag IDs are required")
	}

	// Validate all tag IDs
	for _, id := range tagId {
		if id == uuid.Nil {
			return fmt.Errorf("all tag IDs must be greater than 0")
		}
	}

	// Assign tags to article with context
	err := a.articleTagRepo.AssignTags(ctx, articleId, tagId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to assign tags to article: %v", err)
	}

	return nil
}

// FindArticleByTagId implements ArticleTagService.
func (a *articleTagService) FindArticleByTagId(ctx context.Context, tagId uuid.UUID) ([]model.Article, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate tag ID
	if tagId == uuid.Nil {
		return nil, fmt.Errorf("tag ID must be greater than 0")
	}

	// Get articles by tag from repository with context
	articles, err := a.articleTagRepo.GetArticleByTagId(ctx, tagId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch articles by tag: %v", err)
	}

	return articles, nil
}

// FindTagByArticleId implements ArticleTagService.
func (a *articleTagService) FindTagByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Tags, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate article ID
	if articleId == uuid.Nil {
		return nil, fmt.Errorf("article ID must be greater than 0")
	}

	// Get tags by article from repository with context
	tags, err := a.articleTagRepo.GetTagsByArticleId(ctx, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch tags by article: %v", err)
	}

	return tags, nil
}

func NewArticleTagService(tagRepo repository.TagRepository, articleTagRepo repository.ArticleTagRepository, validationService ValidationService) ArticleTagService {
	return &articleTagService{
		tagRepo:           tagRepo,
		articleTagRepo:    articleTagRepo,
		validationService: validationService,
	}
}
