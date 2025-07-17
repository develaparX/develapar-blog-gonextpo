package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"develapar-server/model/dto"
	"time"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, payload model.Comment) (model.Comment, error)
	GetCommentByArticleId(ctx context.Context, articleId int) ([]model.Comment, error)
	GetCommentByUserId(ctx context.Context, userId int) ([]dto.CommentResponse, error)
	GetCommentById(ctx context.Context, commentId int) (model.Comment, error)
	UpdateComment(ctx context.Context, commentId int, content string, userId int) error
	DeleteComment(ctx context.Context, commentId int) error
}

type commentRepository struct {
	db *sql.DB
}

// GetCommentById implements CommentRepository.
func (c *commentRepository) GetCommentById(ctx context.Context, commentId int) (model.Comment, error) {
	var comment model.Comment
	query := `SELECT id, article_id, user_id, content, created_at FROM comments WHERE id = $1`

	err := c.db.QueryRowContext(ctx, query, commentId).Scan(&comment.Id, &comment.Article.Id, &comment.User.Id, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

// DeleteComment implements CommentRepository.
func (c *commentRepository) DeleteComment(ctx context.Context, commentId int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := c.db.ExecContext(ctx, query, commentId)
	return err
}

// UpdateComment implements CommentRepository.
func (c *commentRepository) UpdateComment(ctx context.Context, commentId int, content string, userId int) error {
	query := `UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 AND user_id=$3`
	_, err := c.db.ExecContext(ctx, query, content, commentId, userId)
	return err
}

// CreateComment implements CommentRepository.
func (c *commentRepository) CreateComment(ctx context.Context, payload model.Comment) (model.Comment, error) {
	var comment model.Comment
	err := c.db.QueryRowContext(ctx, `INSERT INTO comments (article_id, user_id , content, created_at) VALUES($1, $2, $3, $4) RETURNING id, article_id, user_id, content, created_at`, payload.Article.Id, payload.User.Id, payload.Content, time.Now()).Scan(
		&comment.Id, &comment.Article.Id, &comment.User.Id, &comment.Content, &comment.CreatedAt,
	)

	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

// GetCommentByArticleId implements CommentRepository.
func (c *commentRepository) GetCommentByArticleId(ctx context.Context, articleId int) ([]model.Comment, error) {
	var comments []model.Comment

	query := `
	SELECT
		c.id, c.article_id, c.user_id, c.content, c.created_at,
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.created_at, u.updated_at,
		ca.id, ca.name
	FROM comments c
	JOIN articles a ON c.article_id = a.id
	JOIN users u ON c.user_id = u.id
	JOIN categories ca ON a.category_id = ca.id
	WHERE c.article_id = $1
	ORDER BY c.created_at DESC
	`

	rows, err := c.db.QueryContext(ctx, query, articleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment model.Comment
		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&comment.Id, &comment.Article.Id, &comment.User.Id, &comment.Content, &comment.CreatedAt, &article.Id, &article.Title, &article.Slug, &article.Content, &article.Views, &article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, err
		}

		comment.Article = article
		comment.User = user

		comments = append(comments, comment)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetCommentByUserId implements CommentRepository.
func (c *commentRepository) GetCommentByUserId(ctx context.Context, userId int) ([]dto.CommentResponse, error) {
	var comments []dto.CommentResponse

	query := `
	SELECT 
		c.id, c.content, c.created_at,
		u.id, u.name, u.email,
		a.id, a.title, a.slug
	FROM comments c
	JOIN users u ON c.user_id = u.id
	JOIN articles a ON c.article_id = a.id
	WHERE c.user_id = $1
	`

	rows, err := c.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment dto.CommentResponse

		err := rows.Scan(
			&comment.Id, &comment.Content, &comment.CreatedAt,
			&comment.User.Id, &comment.User.Name, &comment.User.Email,
			&comment.Article.Id, &comment.Article.Title, &comment.Article.Slug,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func NewCommentRepository(database *sql.DB) CommentRepository {
	return &commentRepository{db: database}
}
