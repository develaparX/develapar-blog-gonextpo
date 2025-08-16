package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type TagService interface {
	CreateTag(ctx context.Context, payload model.Tags) (model.Tags, error)
	FindAll(ctx context.Context) ([]model.Tags, error)
	FindById(ctx context.Context, id uuid.UUID) (model.Tags, error)
	UpdateTag(ctx context.Context, id uuid.UUID, payload model.Tags) (model.Tags, error)
	DeleteTag(ctx context.Context, id uuid.UUID) error
}

type tagService struct {
	repo              repository.TagRepository
	validationService ValidationService
}

// CreateTag implements TagService.
func (t *tagService) CreateTag(ctx context.Context, payload model.Tags) (model.Tags, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Tags{}, ctx.Err()
	default:
	}

	// Validate tag data
	if strings.TrimSpace(payload.Name) == "" {
		return model.Tags{}, fmt.Errorf("tag name is required")
	}

	// Normalize tag name
	payload.Name = strings.ToLower(strings.TrimSpace(payload.Name))

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.Tags{}, ctx.Err()
	default:
	}

	// Create tag in repository with context
	createdTag, err := t.repo.CreateTag(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Tags{}, ctx.Err()
		}
		return model.Tags{}, fmt.Errorf("failed to create tag: %v", err)
	}

	return createdTag, nil
}

// FindAll implements TagService.
func (t *tagService) FindAll(ctx context.Context) ([]model.Tags, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Get all tags from repository with context
	tags, err := t.repo.GetAllTag(ctx)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch tags: %v", err)
	}

	return tags, nil
}

// FindById implements TagService.
func (t *tagService) FindById(ctx context.Context, id uuid.UUID) (model.Tags, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Tags{}, ctx.Err()
	default:
	}

	// Validate tag ID
	if id == uuid.Nil {
		return model.Tags{}, fmt.Errorf("tag ID must be greater than 0")
	}

	// Get tag by ID from repository with context
	tag, err := t.repo.GetTagById(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Tags{}, ctx.Err()
		}
		return model.Tags{}, fmt.Errorf("failed to fetch tag: %v", err)
	}

	return tag, nil
}

// UpdateTag implements TagService.
func (t *tagService) UpdateTag(ctx context.Context, id uuid.UUID, payload model.Tags) (model.Tags, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Tags{}, ctx.Err()
	default:
	}

	// Validate ID
	if id == uuid.Nil {
		return model.Tags{}, fmt.Errorf("tag ID must be greater than 0")
	}

	// Validate tag data
	if strings.TrimSpace(payload.Name) == "" {
		return model.Tags{}, fmt.Errorf("tag name is required")
	}

	// Normalize tag name
	payload.Name = strings.ToLower(strings.TrimSpace(payload.Name))
	payload.Id = id

	// Check context cancellation before update
	select {
	case <-ctx.Done():
		return model.Tags{}, ctx.Err()
	default:
	}

	// Update tag in repository with context
	updatedTag, err := t.repo.UpdateTag(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Tags{}, ctx.Err()
		}
		return model.Tags{}, fmt.Errorf("failed to update tag: %v", err)
	}

	return updatedTag, nil
}

// DeleteTag implements TagService.
func (t *tagService) DeleteTag(ctx context.Context, id uuid.UUID) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate ID
	if id == uuid.Nil {
		return fmt.Errorf("tag ID must be greater than 0")
	}

	// Delete tag from repository with context
	err := t.repo.DeleteTag(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to delete tag: %v", err)
	}

	return nil
}

func NewTagService(repository repository.TagRepository, validationService ValidationService) TagService {
	return &tagService{
		repo:              repository,
		validationService: validationService,
	}
}
