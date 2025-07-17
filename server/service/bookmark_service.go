package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
	"fmt"
	"strconv"
)

type BookmarkService interface {
	CreateBookmark(ctx context.Context, payload model.Bookmark) (model.Bookmark, error)
	FindByUserId(ctx context.Context, userId string) ([]model.Bookmark, error)
	DeleteBookmark(ctx context.Context, userId, articleId int) error
	IsBookmarked(ctx context.Context, userId, articleId int) (bool, error)
}

type bookmarkService struct {
	repo              repository.BookmarkRepository
	validationService ValidationService
}

// IsBookmarked implements BookmarkService.
func (b *bookmarkService) IsBookmarked(ctx context.Context, userId, articleId int) (bool, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Validate IDs
	if userId <= 0 {
		return false, fmt.Errorf("user ID must be greater than 0")
	}
	if articleId <= 0 {
		return false, fmt.Errorf("article ID must be greater than 0")
	}

	// Check bookmark status with context
	isBookmarked, err := b.repo.IsBookmarked(ctx, userId, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		return false, fmt.Errorf("failed to check bookmark status: %v", err)
	}

	return isBookmarked, nil
}

// DeleteBookmark implements BookmarkService.
func (b *bookmarkService) DeleteBookmark(ctx context.Context, userId, articleId int) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate IDs
	if userId <= 0 {
		return fmt.Errorf("user ID must be greater than 0")
	}
	if articleId <= 0 {
		return fmt.Errorf("article ID must be greater than 0")
	}

	// Delete bookmark from repository with context
	err := b.repo.DeleteBookmark(ctx, userId, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to delete bookmark: %v", err)
	}

	return nil
}

// CreateBookmark implements BookmarkService.
func (b *bookmarkService) CreateBookmark(ctx context.Context, payload model.Bookmark) (model.Bookmark, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Bookmark{}, ctx.Err()
	default:
	}

	// Validate bookmark data
	if payload.User.Id <= 0 {
		return model.Bookmark{}, fmt.Errorf("valid user ID is required")
	}
	if payload.Article.Id <= 0 {
		return model.Bookmark{}, fmt.Errorf("valid article ID is required")
	}

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.Bookmark{}, ctx.Err()
	default:
	}

	// Create bookmark in repository with context
	createdBookmark, err := b.repo.CreateBookmark(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Bookmark{}, ctx.Err()
		}
		return model.Bookmark{}, fmt.Errorf("failed to create bookmark: %v", err)
	}

	return createdBookmark, nil
}

// FindByUserId implements BookmarkService.
func (b *bookmarkService) FindByUserId(ctx context.Context, userId string) ([]model.Bookmark, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate user ID
	if userId == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// Convert string ID to int for validation
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}
	if userIdInt <= 0 {
		return nil, fmt.Errorf("user ID must be greater than 0")
	}

	// Get bookmarks by user from repository with context
	bookmarks, err := b.repo.GetByUserId(ctx, userId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch bookmarks by user: %v", err)
	}

	return bookmarks, nil
}

func NewBookmarkService(repository repository.BookmarkRepository, validationService ValidationService) BookmarkService {
	return &bookmarkService{
		repo:              repository,
		validationService: validationService,
	}
}
