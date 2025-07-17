package service

import (
	"context"
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
	return b.repo.IsBookmarked(context.Background(), userId, articleId)
}

// DeleteBookmark implements BookmarkService.
func (b *bookmarkService) DeleteBookmark(userId int, articleId int) error {
	return b.repo.DeleteBookmark(context.Background(), userId, articleId)
}

// CreateBookmark implements BookmarkService.
func (b *bookmarkService) CreateBookmark(payload model.Bookmark) (model.Bookmark, error) {
	return b.repo.CreateBookmark(context.Background(), payload)
}

// FindByUserId implements BookmarkService.
func (b *bookmarkService) FindByUserId(userId string) ([]model.Bookmark, error) {
	return b.repo.GetByUserId(context.Background(), userId)
}

func NewBookmarkService(repository repository.BookmarkRepository) BookmarkService {
	return &bookmarkService{repo: repository}
}
