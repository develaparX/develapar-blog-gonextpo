package repository

import (
	"database/sql"
	"develapar-server/model"
	"time"
)

type LikeRepository interface {
	CreateLike(payload model.Likes) (model.Likes, error)
	GetLikeByArticleId(articleId int) ([]model.Likes, error)
	GetLikeByUserId(userId int) ([]model.Likes, error)
	DeleteLike(userId, articleId int) error
	IsLiked(userId, articleId int) (bool, error)
}

type likeRepository struct {
	db *sql.DB
}

// isLiked implements LikeRepository.
func (r *likeRepository) IsLiked(userId int, articleId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND article_id = $2)`
	err := r.db.QueryRow(query, userId, articleId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateLike implements LikeRepository.
func (l *likeRepository) CreateLike(payload model.Likes) (model.Likes, error) {
	var like model.Likes
	err := l.db.QueryRow(`INSERT INTO likes(article_id, user_id, created_at) VALUES($1,$2,$3) RETURNING id, article_id, user_id, created_at `, payload.Article.Id, payload.User.Id, time.Now()).Scan(
		&like.Id, &like.Article.Id, &like.User.Id, &like.CreatedAt,
	)

	if err != nil {
		return model.Likes{}, err
	}

	return like, nil
}

// DeleteLike implements LikeRepository.
func (l *likeRepository) DeleteLike(userId int, articleId int) error {
	query := `DELETE FROM likes WHERE user_id=$1 AND article_id = $2`
	_, err := l.db.Exec(query, userId, articleId)

	return err
}

// GetLikeByArticleId implements LikeRepository.
func (l *likeRepository) GetLikeByArticleId(articleId int) ([]model.Likes, error) {
	var likes []model.Likes

	query := `
	SELECT
		l.id, l.article_id, l.user_id, l.created_at,
		u.id, u.name, u.email
	FROM likes l
	JOIN users u ON l.user_id = u.id
	WHERE l.article_id = $1

	`

	rows, err := l.db.Query(query, articleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var like model.Likes
		var user model.User

		err := rows.Scan(
			&like.Id, &like.Article.Id, &like.User.Id, &like.CreatedAt, &user.Id, &user.Name, &user.Email,
		)

		if err != nil {
			return nil, err
		}

		like.User = user

		likes = append(likes, like)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil

}

// GetLikeByUserId implements LikeRepository.
func (l *likeRepository) GetLikeByUserId(userId int) ([]model.Likes, error) {
	var likes []model.Likes

	query := `
	SELECT
		l.id, l.article_id, l.user_id, l.created_at,
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at
	FROM likes l
	JOIN articles a ON l.article_id = a.id
	WHERE l.user_id = $1

	`

	rows, err := l.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var like model.Likes
		var article model.Article

		err := rows.Scan(
			&like.Id, &like.Article.Id, &like.User.Id, &like.CreatedAt, &article.Id, &article.Title, &article.Slug, &article.Content, &article.Views, &article.CreatedAt, &article.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		like.Article = article

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
