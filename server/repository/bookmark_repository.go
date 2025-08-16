package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"time"

	"github.com/google/uuid"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, payload model.Bookmark) (model.Bookmark, error)
	GetByUserId(ctx context.Context, userId uuid.UUID) ([]model.Bookmark, error)
	DeleteBookmark(ctx context.Context, userId, articleId uuid.UUID) error
	IsBookmarked(ctx context.Context, userId, articleId uuid.UUID) (bool, error)
}

type bookmarkRepository struct {
	db *sql.DB
}

// DeleteBookmark implements BookmarkRepository.
func (b *bookmarkRepository) DeleteBookmark(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) error {
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2`
	_, err := b.db.ExecContext(ctx, query, userId, articleId)
	return err
}

// IsBookmarked implements BookmarkRepository.
func (r *bookmarkRepository) IsBookmarked(ctx context.Context, userId uuid.UUID, articleId uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND article_id = $2)`
	err := r.db.QueryRowContext(ctx, query, userId, articleId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetByUserId implements BookmarkRepository.
func (b *bookmarkRepository) GetByUserId(ctx context.Context, userId uuid.UUID) ([]model.Bookmark, error) {
	var bookmarks []model.Bookmark

	query := `
	SELECT 
		b.id, b.article_id, b.user_id, b.created_at, b.updated_at,
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role, u.created_at, u.updated_at,
		c.id, c.name, c.created_at, c.updated_at
	FROM bookmarks b
	JOIN articles a ON b.article_id = a.id
	JOIN users u ON b.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE b.user_id = $1
	ORDER BY b.created_at DESC;
	`

	rows, err := b.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var bookmark model.Bookmark
		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&bookmark.Id, &bookmark.ArticleId, &bookmark.UserId, &bookmark.CreatedAt, &bookmark.UpdatedAt,
			&article.Id, &article.Title, &article.Slug, &article.Content, &article.UserId, &article.CategoryId, &article.Views, &article.Status, &article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt,
			&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Assigning the article, user, and category to the bookmark
		bookmark.Article = &article
		bookmark.User = &user
		article.User = &user
		article.Category = &category

		// Append the bookmark to the bookmarks slice
		bookmarks = append(bookmarks, bookmark)
	}

	// Handle errors while iterating through rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookmarks, nil
}

// CreateBookmark implements BookmarkRepository.
func (b *bookmarkRepository) CreateBookmark(ctx context.Context, payload model.Bookmark) (model.Bookmark, error) {
	newId := uuid.Must(uuid.NewV7())
	var brk model.Bookmark
	err := b.db.QueryRowContext(ctx, `INSERT INTO bookmarks (id, article_id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, article_id, user_id, created_at, updated_at`, newId, payload.ArticleId, payload.UserId, time.Now(), time.Now()).Scan(
		&brk.Id, &brk.ArticleId, &brk.UserId, &brk.CreatedAt, &brk.UpdatedAt,
	)

	if err != nil {
		return model.Bookmark{}, err
	}
	return brk, nil
}

func NewBookmarkRepository(database *sql.DB) BookmarkRepository {
	return &bookmarkRepository{db: database}
}
