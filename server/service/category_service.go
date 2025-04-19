package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type CategoryService interface {
	CreateCategory(payload model.Category) (model.Category, error)
	FindAll() ([]model.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

// CreateCategory implements CategoryService.
func (c *categoryService) CreateCategory(payload model.Category) (model.Category, error) {
	return c.repo.CreateCategory(payload)
}

// FindAll implements CategoryService.
func (c *categoryService) FindAll() ([]model.Category, error) {
	return c.repo.GetAll()
}

func NewCategoryService(repository repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repository}
}
