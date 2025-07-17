package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"errors"
)

type CommentService interface {
	CreateComment(payload model.Comment) (model.Comment, error)
	FindCommentByArticleId(articleId int) ([]model.Comment, error)
	FindCommentByUserId(userId int) ([]dto.CommentResponse, error)
	EditComment(commentId int, content string, userId int) error
	DeleteComment(commentId int, userId int) error
}

var ErrUnauthorized = errors.New("unauthorized")

type commentService struct {
	repo repository.CommentRepository
}

// DeleteComment implements CommentService.
func (c *commentService) DeleteComment(commentId int, userId int) error {
	comment, err := c.repo.GetCommentById(context.Background(), commentId)
	if err != nil {
		return err
	}

	if comment.User.Id != userId {
		return ErrUnauthorized
	}

	return c.repo.DeleteComment(context.Background(), commentId)
}

// EditComment implements CommentService.
func (c *commentService) EditComment(commentId int, content string, userId int) error {
	comment, err := c.repo.GetCommentById(context.Background(), commentId)
	if err != nil {
		return err
	}

	if comment.User.Id != userId {
		return ErrUnauthorized
	}

	return c.repo.UpdateComment(context.Background(), commentId, content, userId)
}

// CreateComment implements CommentService.
func (c *commentService) CreateComment(payload model.Comment) (model.Comment, error) {
	return c.repo.CreateComment(context.Background(), payload)
}

// FindCommentByArticleId implements CommentService.
func (c *commentService) FindCommentByArticleId(articleId int) ([]model.Comment, error) {
	return c.repo.GetCommentByArticleId(context.Background(), articleId)
}

// FindCommentByUserId implements CommentService.
func (c *commentService) FindCommentByUserId(userId int) ([]dto.CommentResponse, error) {
	return c.repo.GetCommentByUserId(context.Background(), userId)
}

func NewCommentService(repository repository.CommentRepository) CommentService {
	return &commentService{repo: repository}
}
