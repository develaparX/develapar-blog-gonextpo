package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"develapar-server/model/dto"
	"time"

	"github.com/google/uuid"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, payload model.Comment) (model.Comment, error)
	GetCommentByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Comment, error)
	GetCommentByUserId(ctx context.Context, userId uuid.UUID) ([]dto.CommentResponse, error)
	GetCommentById(ctx context.Context, commentId uuid.UUID) (model.Comment, error)
	UpdateComment(ctx context.Context, commentId uuid.UUID, content string, userId uuid.UUID) error
	DeleteComment(ctx context.Context, commentId uuid.UUID) error
}

type commentRepository struct {
	db *sql.DB
}

// GetCommentById implements CommentRepository.
func (c *commentRepository) GetCommentById(ctx context.Context, commentId uuid.UUID) (model.Comment, error) {
	var comment model.Comment
	query := `SELECT id, article_id, user_id, content, created_at, updated_at FROM comments WHERE id = $1`

	err := c.db.QueryRowContext(ctx, query, commentId).Scan(&comment.Id, &comment.ArticleId, &comment.UserId, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

// DeleteComment implements CommentRepository.
func (c *commentRepository) DeleteComment(ctx context.Context, commentId uuid.UUID) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := c.db.ExecContext(ctx, query, commentId)
	return err
}

// UpdateComment implements CommentRepository.
func (c *commentRepository) UpdateComment(ctx context.Context, commentId uuid.UUID, content string, userId uuid.UUID) error {
	query := `UPDATE comments SET content = $1, updated_at = NOW() WHERE id = $2 AND user_id=$3`
	_, err := c.db.ExecContext(ctx, query, content, commentId, userId)
	return err
}

// CreateComment implements CommentRepository.
func (c *commentRepository) CreateComment(ctx context.Context, payload model.Comment) (model.Comment, error) {
	newId := uuid.Must(uuid.NewV7())
	var comment model.Comment
	err := c.db.QueryRowContext(ctx, `INSERT INTO comments (id, article_id, user_id, content, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, article_id, user_id, content, created_at, updated_at`, newId, payload.ArticleId, payload.UserId, payload.Content, time.Now(), time.Now()).Scan(
		&comment.Id, &comment.ArticleId, &comment.UserId, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
	)

	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

// GetCommentByArticleId implements CommentRepository.
func (c *commentRepository) GetCommentByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Comment, error) {
	var comments []model.Comment

	query := `
	SELECT
		c.id, c.article_id, c.user_id, c.content, c.created_at, c.updated_at,
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role, u.created_at, u.updated_at,
		ca.id, ca.name, ca.created_at, ca.updated_at
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
			&comment.Id, &comment.ArticleId, &comment.UserId, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&article.Id, &article.Title, &article.Slug, &article.Content, &article.UserId, &article.CategoryId, &article.Views, &article.Status, &article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
			&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		comment.Article = &article
		comment.User = &user
		article.User = &user
		article.Category = &category

		comments = append(comments, comment)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetCommentByUserId implements CommentRepository.
func (c *commentRepository) GetCommentByUserId(ctx context.Context, userId uuid.UUID) ([]dto.CommentResponse, error) {
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
