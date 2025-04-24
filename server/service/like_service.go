package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type LikeService interface {
	CreateLike(payload model.Likes) (model.Likes, error)
	FindLikeByArticleId(articleId int) ([]model.Likes, error)
	FindLikeByUserId(userId int) ([]model.Likes, error)
	DeleteLike(userId, articleId int) error
}

type likeService struct {
	repo repository.LikeRepository
}

// CreateLike implements LikeService.
func (l *likeService) CreateLike(payload model.Likes) (model.Likes, error) {
	return l.repo.CreateLike(payload)
}

// DeleteLike implements LikeService.
func (l *likeService) DeleteLike(userId int, articleId int) error {
	return l.repo.DeleteLike(userId, articleId)
}

// FindLikeByArticleId implements LikeService.
func (l *likeService) FindLikeByArticleId(articleId int) ([]model.Likes, error) {

	return l.repo.GetLikeByArticleId(articleId)
}

// FindLikeByUserId implements LikeService.
func (l *likeService) FindLikeByUserId(userId int) ([]model.Likes, error) {
	return l.repo.GetLikeByUserId(userId)
}

func NewLikeService(repository repository.LikeRepository) LikeService {
	return &likeService{repo: repository}
}
