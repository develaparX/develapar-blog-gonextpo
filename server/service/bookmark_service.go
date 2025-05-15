package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type BookmarkService interface {
	CreateBookmark(payload model.Bookmark) (model.Bookmark, error)
	FindByUserId(userId string) ([]model.Bookmark, error)
	DeleteBookmark(userId, articleId int) error
	IsBookmarked(userId, articleId int) (bool, error)
}

type bookmarkService struct {
	repo repository.BookmarkRepository
}

// IsBookmarked implements BookmarkService.
func (b *bookmarkService) IsBookmarked(userId int, articleId int) (bool, error) {
	return b.repo.IsBookmarked(userId, articleId)
}

// DeleteBookmark implements BookmarkService.
func (b *bookmarkService) DeleteBookmark(userId int, articleId int) error {
	return b.repo.DeleteBookmark(userId, articleId)
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
