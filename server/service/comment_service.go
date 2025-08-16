package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"errors"

	"github.com/google/uuid"
)

type CommentService interface {
	CreateComment(ctx context.Context, payload model.Comment) (model.Comment, error)
	FindCommentByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Comment, error)
	FindCommentByUserId(ctx context.Context, userId uuid.UUID) ([]dto.CommentResponse, error)
	EditComment(ctx context.Context, commentId uuid.UUID, content string, userId uuid.UUID) error
	DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) error
}

var ErrUnauthorized = errors.New("unauthorized")

type commentService struct {
	repo              repository.CommentRepository
	validationService ValidationService
}

// DeleteComment implements CommentService.
func (c *commentService) DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate IDs
	if commentId == uuid.Nil {
		return errors.New("comment ID must be greater than 0")
	}
	if userId == uuid.Nil {
		return errors.New("user ID must be greater than 0")
	}

	// Get comment with context
	comment, err := c.repo.GetCommentById(ctx, commentId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	// Check authorization
	if comment.User.Id != userId {
		return ErrUnauthorized
	}

	// Check context cancellation before delete
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Delete comment with context
	err = c.repo.DeleteComment(ctx, commentId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	return nil
}

// EditComment implements CommentService.
func (c *commentService) EditComment(ctx context.Context, commentId uuid.UUID, content string, userId uuid.UUID) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate IDs and content
	if commentId == uuid.Nil {
		return errors.New("comment ID must be greater than 0")
	}
	if userId == uuid.Nil {
		return errors.New("user ID must be greater than 0")
	}
	if content == "" {
		return errors.New("comment content is required")
	}

	// Get comment with context
	comment, err := c.repo.GetCommentById(ctx, commentId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	// Check authorization
	if comment.User.Id != userId {
		return ErrUnauthorized
	}

	// Create comment object for validation
	updatedComment := comment
	updatedComment.Content = content

	// Validate comment data using validation service
	if validationErr := c.validationService.ValidateComment(ctx, updatedComment); validationErr != nil {
		return validationErr
	}

	// Check context cancellation before update
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Update comment with context
	err = c.repo.UpdateComment(ctx, commentId, content, userId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	return nil
}

// CreateComment implements CommentService.
func (c *commentService) CreateComment(ctx context.Context, payload model.Comment) (model.Comment, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Comment{}, ctx.Err()
	default:
	}

	// Validate comment data using validation service
	if validationErr := c.validationService.ValidateComment(ctx, payload); validationErr != nil {
		return model.Comment{}, validationErr
	}

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.Comment{}, ctx.Err()
	default:
	}

	// Create comment in repository with context
	createdComment, err := c.repo.CreateComment(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Comment{}, ctx.Err()
		}
		return model.Comment{}, err
	}

	return createdComment, nil
}

// FindCommentByArticleId implements CommentService.
func (c *commentService) FindCommentByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Comment, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate article ID
	if articleId == uuid.Nil {
		return nil, errors.New("article ID must be greater than 0")
	}

	// Get comments by article from repository with context
	comments, err := c.repo.GetCommentByArticleId(ctx, articleId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}

	return comments, nil
}

// FindCommentByUserId implements CommentService.
func (c *commentService) FindCommentByUserId(ctx context.Context, userId uuid.UUID) ([]dto.CommentResponse, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate user ID
	if userId == uuid.Nil {
		return nil, errors.New("user ID must be greater than 0")
	}

	// Get comments by user from repository with context
	comments, err := c.repo.GetCommentByUserId(ctx, userId)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}

	return comments, nil
}

func NewCommentService(repository repository.CommentRepository, validationService ValidationService) CommentService {
	return &commentService{
		repo:              repository,
		validationService: validationService,
	}
}
