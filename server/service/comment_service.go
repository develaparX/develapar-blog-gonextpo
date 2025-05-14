package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
)

type CommentService interface {
	CreateComment(payload model.Comment) (model.Comment, error)
	FindCommentByArticleId(articleId int) ([]model.Comment, error)
	FindCommentByUserId(userId int) ([]dto.CommentResponse, error)
	EditComment(commentId int, content string) error
	DeleteComment(commentId int) error
}

type commentService struct {
	repo repository.CommentRepository
}

// DeleteComment implements CommentService.
func (c *commentService) DeleteComment(commentId int) error {
	return c.repo.DeleteComment(commentId)
}

// EditComment implements CommentService.
func (c *commentService) EditComment(commentId int, content string) error {
	return c.repo.UpdateComment(commentId, content)
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
func (c *commentService) FindCommentByUserId(userId int) ([]dto.CommentResponse, error) {
	return c.repo.GetCommentByUserId(userId)
}

func NewCommentService(repository repository.CommentRepository) CommentService {
	return &commentService{repo: repository}
}
