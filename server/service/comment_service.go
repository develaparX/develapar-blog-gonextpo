package service

import (
	"develapar-server/model"
	"develapar-server/repository"
)

type CommentService interface {
	CreateComment(payload model.Comment) (model.Comment, error)
	FindCommentByArticleId(articleId int) ([]model.Comment, error)
	FindCommentByUserId(userId int) ([]model.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

// CreateComment implements CommentService.
func (c *commentService) CreateComment(payload model.Comment) (model.Comment, error) {
	return c.repo.CreateComment(payload)
}

// FindCommentByArticleId implements CommentService.
func (c *commentService) FindCommentByArticleId(articleId int) ([]model.Comment, error) {
	return c.repo.GetCommentByArticleId(articleId)
}

// FindCommentByUserId implements CommentService.
func (c *commentService) FindCommentByUserId(userId int) ([]model.Comment, error) {
	return c.repo.GetCommentByUserId(userId)
}

func NewCommentService(repository repository.CommentRepository) CommentService {
	return &commentService{repo: repository}
}
