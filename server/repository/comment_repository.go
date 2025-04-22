package repository

import (
	"database/sql"
	"develapar-server/model"
	"time"
)

type CommentRepository interface {
	CreateComment(payload model.Comment) (model.Comment, error)
	GetCommentByArticleId(articleId int) ([]model.Comment, error)
	GetCommentByUserId(userId int) ([]model.Comment, error)
}

type commentRepository struct {
	db *sql.DB
}

// CreateComment implements CommentRepository.
func (c *commentRepository) CreateComment(payload model.Comment) (model.Comment, error) {
	var comment model.Comment
	err := c.db.QueryRow(`INSERT INTO comments (article_id, user_id , content, created_at) VALUES($1, $2, $3, $4) RETURNING id, article_id, user_id, content, created_at`, payload.Article.Id, payload.User.Id, payload.Content, time.Now()).Scan(
		&comment.Id, &comment.Article.Id, &comment.User.Id, &comment.Content, &comment.CreatedAt,
	)

	if err != nil {
		return model.Comment{}, err
	}

	return comment, nil
}

// GetCommentByArticleId implements CommentRepository.
func (c *commentRepository) GetCommentByArticleId(articleId int) ([]model.Comment, error) {
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

	rows, err := c.db.Query(query, articleId)
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
func (c *commentRepository) GetCommentByUserId(userId int) ([]model.Comment, error) {
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
	WHERE c.user_id = $1
	ORDER BY c.created_at DESC
	`

	rows, err := c.db.Query(query, userId)
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

func NewCommentRepository(database *sql.DB) CommentRepository {
	return &commentRepository{db: database}
}
