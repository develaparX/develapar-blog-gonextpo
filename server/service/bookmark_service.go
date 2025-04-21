package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type BookmarkService interface {
	CreateBookmark(payload model.Bookmark) (model.Bookmark, error)
	FindByUserId(userId string) ([]model.Bookmark, error)
}

type bookmarkService struct {
	repo repository.BookmarkRepository
}

// CreateBookmark implements BookmarkService.
func (b *bookmarkService) CreateBookmark(payload model.Bookmark) (model.Bookmark, error) {
	return b.repo.CreateBookmark(payload)
}

// FindByUserId implements BookmarkService.
func (b *bookmarkService) FindByUserId(userId string) ([]model.Bookmark, error) {
	return b.repo.GetByUserId(userId)
}

func NewBookmarkService(repository repository.BookmarkRepository) BookmarkService {
	return &bookmarkService{repo: repository}
}
