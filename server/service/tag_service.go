package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
)

type TagService interface {
	CreateTag(payload model.Tags) (model.Tags, error)
	FindAll() ([]model.Tags, error)
	FindById(id int) (model.Tags, error)
}

type tagService struct {
	repo repository.TagRepository
}

// CreateTag implements TagService.
func (t *tagService) CreateTag(payload model.Tags) (model.Tags, error) {
	return t.repo.CreateTag(context.Background(), payload)
}

// FindAll implements TagService.
func (t *tagService) FindAll() ([]model.Tags, error) {
	return t.repo.GetAllTag(context.Background())
}

// FindById implements TagService.
func (t *tagService) FindById(id int) (model.Tags, error) {

	return t.repo.GetTagById(context.Background(), id)
}

func NewTagService(repository repository.TagRepository) TagService {
	return &tagService{repo: repository}
}
