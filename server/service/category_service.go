package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"strings"
)

type CategoryService interface {
	CreateCategory(payload model.Category) (model.Category, error)
	FindAll() ([]model.Category, error)
	UpdateCategory(id int, req dto.UpdateCategoryRequest) (model.Category, error)
	DeleteCategory(id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

// DeleteCategory implements CategoryService.
func (c *categoryService) DeleteCategory(id int) error {
	return c.repo.DeleteCategory(id)
}

// UpdateCategory implements CategoryService.
func (c *categoryService) UpdateCategory(id int, req dto.UpdateCategoryRequest) (model.Category, error) {
	cat, err := c.repo.GetCategoryById(id)
	if err != nil {
		return model.Category{}, err
	}

	if req.Name != nil {
		cat.Name = strings.ToLower(*req.Name)
	}

	return c.repo.UpdateCategory(cat)

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
