package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"time"

	"github.com/google/uuid"
)

type LikeRepository interface {
	CreateLike(ctx context.Context, payload model.Likes) (model.Likes, error)
	GetLikeByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Likes, error)
	GetLikeByUserId(ctx context.Context, userId uuid.UUID) ([]model.Likes, error)
	DeleteLike(ctx context.Context, userId, articleId uuid.UUID) error
	IsLiked(ctx context.Context, userId, articleId uuid.UUID) (bool, error)
}

type likeRepository struct {
	db *sql.DB
}

// IsLiked implements LikeRepository.
func (r *likeRepository) IsLiked(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND article_id = $2)`
	err := r.db.QueryRowContext(ctx, query, userId, articleId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateLike implements LikeRepository.
func (l *likeRepository) CreateLike(ctx context.Context, payload model.Likes) (model.Likes, error) {
	newId := uuid.Must(uuid.NewV7())
	var like model.Likes
	err := l.db.QueryRowContext(ctx, `INSERT INTO likes(id, article_id, user_id, created_at, updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id, article_id, user_id, created_at, updated_at`, newId, payload.ArticleId, payload.UserId, time.Now(), time.Now()).Scan(
		&like.Id, &like.ArticleId, &like.UserId, &like.CreatedAt, &like.UpdatedAt,
	)

	if err != nil {
		return model.Likes{}, err
	}

	return like, nil
}

// DeleteLike implements LikeRepository.
func (l *likeRepository) DeleteLike(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) error {
	query := `DELETE FROM likes WHERE user_id=$1 AND article_id = $2`
	_, err := l.db.ExecContext(ctx, query, userId, articleId)

	return err
}

// GetLikeByArticleId implements LikeRepository.
func (l *likeRepository) GetLikeByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Likes, error) {
	var likes []model.Likes

	query := `
	SELECT
		l.id, l.article_id, l.user_id, l.created_at, l.updated_at,
		u.id, u.name, u.email, u.role
	FROM likes l
	JOIN users u ON l.user_id = u.id
	WHERE l.article_id = $1

	`

	rows, err := l.db.QueryContext(ctx, query, articleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var like model.Likes
		var user model.User

		err := rows.Scan(
			&like.Id, &like.ArticleId, &like.UserId, &like.CreatedAt, &like.UpdatedAt, &user.Id, &user.Name, &user.Email, &user.Role,
		)

		if err != nil {
			return nil, err
		}

		like.User = &user

		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil

}

// GetLikeByUserId implements LikeRepository.
func (l *likeRepository) GetLikeByUserId(ctx context.Context, userId uuid.UUID) ([]model.Likes, error) {
	var likes []model.Likes

	query := `
	SELECT
		l.id, l.article_id, l.user_id, l.created_at, l.updated_at,
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at
	FROM likes l
	JOIN articles a ON l.article_id = a.id
	WHERE l.user_id = $1

	`

	rows, err := l.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var like model.Likes
		var article model.Article

		err := rows.Scan(
			&like.Id, &like.ArticleId, &like.UserId, &like.CreatedAt, &like.UpdatedAt, 
			&article.Id, &article.Title, &article.Slug, &article.Content, &article.UserId, &article.CategoryId, &article.Views, &article.Status, &article.CreatedAt, &article.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		like.Article = &article

		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func NewLikeRepository(database *sql.DB) LikeRepository {
	return &likeRepository{db: database}
}
