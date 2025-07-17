package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"fmt"
	"strings"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, payload model.Category) (model.Category, error)
	FindAll(ctx context.Context) ([]model.Category, error)
	UpdateCategory(ctx context.Context, id int, req dto.UpdateCategoryRequest) (model.Category, error)
	DeleteCategory(ctx context.Context, id int) error
}

type categoryService struct {
	repo              repository.CategoryRepository
	validationService ValidationService
}

// DeleteCategory implements CategoryService.
func (c *categoryService) DeleteCategory(ctx context.Context, id int) error {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate ID
	if id <= 0 {
		return fmt.Errorf("category ID must be greater than 0")
	}

	// Delete category from repository with context
	err := c.repo.DeleteCategory(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("failed to delete category: %v", err)
	}

	return nil
}

// UpdateCategory implements CategoryService.
func (c *categoryService) UpdateCategory(ctx context.Context, id int, req dto.UpdateCategoryRequest) (model.Category, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Category{}, ctx.Err()
	default:
	}

	// Validate ID
	if id <= 0 {
		return model.Category{}, fmt.Errorf("category ID must be greater than 0")
	}

	// Get existing category with context
	cat, err := c.repo.GetCategoryById(ctx, id)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Category{}, ctx.Err()
		}
		return model.Category{}, fmt.Errorf("failed to fetch category for update: %v", err)
	}

	// Update fields if provided
	if req.Name != nil {
		cat.Name = strings.ToLower(*req.Name)
	}

	// Basic validation for category name
	if strings.TrimSpace(cat.Name) == "" {
		return model.Category{}, fmt.Errorf("category name is required")
	}

	// Check context cancellation before update
	select {
	case <-ctx.Done():
		return model.Category{}, ctx.Err()
	default:
	}

	// Update category in repository with context
	updatedCategory, err := c.repo.UpdateCategory(ctx, cat)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Category{}, ctx.Err()
		}
		return model.Category{}, fmt.Errorf("failed to update category: %v", err)
	}

	return updatedCategory, nil
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(ctx context.Context, payload model.Category) (model.Category, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return model.Category{}, ctx.Err()
	default:
	}

	// Basic validation for category
	if strings.TrimSpace(payload.Name) == "" {
		return model.Category{}, fmt.Errorf("category name is required")
	}

	// Normalize category name
	payload.Name = strings.ToLower(strings.TrimSpace(payload.Name))

	// Check context cancellation after validation
	select {
	case <-ctx.Done():
		return model.Category{}, ctx.Err()
	default:
	}

	// Create category in repository with context
	createdCategory, err := c.repo.CreateCategory(ctx, payload)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return model.Category{}, ctx.Err()
		}
		return model.Category{}, fmt.Errorf("failed to create category: %v", err)
	}

	return createdCategory, nil
}

// FindAll implements CategoryService.
func (c *categoryService) FindAll(ctx context.Context) ([]model.Category, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Get all categories from repository with context
	categories, err := c.repo.GetAll(ctx)
	if err != nil {
		// Check if context was cancelled during repository operation
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("failed to fetch categories: %v", err)
	}

	return categories, nil
}

func NewCategoryService(repository repository.CategoryRepository, validationService ValidationService) CategoryService {
	return &categoryService{
		repo:              repository,
		validationService: validationService,
	}
}
