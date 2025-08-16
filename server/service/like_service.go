package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
	"fmt"

	"github.com/google/uuid"
)

type LikeService interface {
	CreateLike(ctx context.Context, payload model.Likes) (model.Likes, error)
	FindLikeByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Likes, error)
	FindLikeByUserId(ctx context.Context, userId uuid.UUID) ([]model.Likes, error)
	DeleteLike(ctx context.Context, userId, articleId uuid.UUID) error
	IsLiked(ctx context.Context, userId, articleId uuid.UUID) (bool, error)
}

type likeService struct {
	repo              repository.LikeRepository
	validationService ValidationService
}

// IsLiked implements LikeService.
func (l *likeService) IsLiked(ctx context.Context, userId, articleId uuid.UUID) (bool, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	// Validate IDs
	if userId == uuid.Nil {
		return false, fmt.Errorf("user ID must be greater than 0")
	}
	if articleId == uuid.Nil {
		return false, fmt.Errorf("article ID must be greater than 0")
	}

	// Check like status with context
	isLiked, err := l.repo.IsLiked(ctx, userId, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return false, ctx.Err()
		}
		return false, fmt.Errorf("failed to check like status: %v", err)
	}

	return isLiked, nil
}

// CreateLike implements LikeService.
func (l *likeService) CreateLike(ctx context.Context, payload model.Likes) (model.Likes, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Likes{}, ctx.Err()
	default:
	}

	// Validate like data
	if payload.User.Id == uuid.Nil {
		return model.Likes{}, fmt.Errorf("valid user ID is required")
	}
	if payload.Article.Id == uuid.Nil {
		return model.Likes{}, fmt.Errorf("valid article ID is required")
	}

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.Likes{}, ctx.Err()
	default:
	}

	// Create like in repository with context
	createdLike, err := l.repo.CreateLike(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Likes{}, ctx.Err()
		}
		return model.Likes{}, fmt.Errorf("failed to create like: %v", err)
	}

	return createdLike, nil
}

// DeleteLike implements LikeService.
func (l *likeService) DeleteLike(ctx context.Context, userId, articleId uuid.UUID) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate IDs
	if userId == uuid.Nil {
		return fmt.Errorf("user ID must be greater than 0")
	}
	if articleId == uuid.Nil {
		return fmt.Errorf("article ID must be greater than 0")
	}

	// Delete like from repository with context
	err := l.repo.DeleteLike(ctx, userId, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to delete like: %v", err)
	}

	return nil
}

// FindLikeByArticleId implements LikeService.
func (l *likeService) FindLikeByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Likes, error) {
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

	// Get likes by article from repository with context
	likes, err := l.repo.GetLikeByArticleId(ctx, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch likes by article: %v", err)
	}

	return likes, nil
}

// FindLikeByUserId implements LikeService.
func (l *likeService) FindLikeByUserId(ctx context.Context, userId uuid.UUID) ([]model.Likes, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate user ID
	if userId == uuid.Nil {
		return nil, fmt.Errorf("user ID must be greater than 0")
	}

	// Get likes by user from repository with context
	likes, err := l.repo.GetLikeByUserId(ctx, userId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch likes by user: %v", err)
	}

	return likes, nil
}

func NewLikeService(repository repository.LikeRepository, validationService ValidationService) LikeService {
	return &likeService{
		repo:              repository,
		validationService: validationService,
	}
}
