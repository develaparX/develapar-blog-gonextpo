package repository

import (
	"database/sql"
	"develapar-server/model"
	"time"
)

type BookmarkRepository interface {
	CreateBookmark(payload model.Bookmark) (model.Bookmark, error)
	GetByUserId(userId string) ([]model.Bookmark, error)
	DeleteBookmark(userId, articleId int) error
	IsBookmarked(userId, articleId int) (bool, error)
}

type bookmarkRepository struct {
	db *sql.DB
}

// DeleteBookmark implements BookmarkRepository.
func (b *bookmarkRepository) DeleteBookmark(userId int, articleId int) error {
	query := `DELETE FROM bookmarks WHERE user_id = $1 AND article_id = $2`
	_, err := b.db.Exec(query, userId, articleId)
	return err
}

// IsBookmarked implements BookmarkRepository.
func (r *bookmarkRepository) IsBookmarked(userId, articleId int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM bookmarks WHERE user_id = $1 AND article_id = $2)`
	err := r.db.QueryRow(query, userId, articleId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetByUserId implements BookmarkRepository.
func (b *bookmarkRepository) GetByUserId(userId string) ([]model.Bookmark, error) {
	var bookmarks []model.Bookmark

	query := `
	SELECT 
		b.id, b.article_id, b.user_id, b.created_at, 
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.created_at, u.updated_at,
		c.id, c.name
	FROM bookmarks b
	JOIN articles a ON b.article_id = a.id
	JOIN users u ON b.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE b.user_id = $1
	ORDER BY b.created_at DESC;
	`

	rows, err := b.db.Query(query, userId)
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
			&bookmark.Id, &article.Id, &bookmark.User.Id, &bookmark.CreatedAt,
			&article.Id, &article.Title, &article.Slug, &article.Content, &article.Views, &article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, err
		}

		// Assigning the article, user, and category to the bookmark
		bookmark.Article = article
		bookmark.User = user

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
func (b *bookmarkRepository) CreateBookmark(payload model.Bookmark) (model.Bookmark, error) {
	var brk model.Bookmark
	err := b.db.QueryRow(`INSERT INTO bookmarks (article_id, user_id, created_at) VALUES ($1, $2, $3) RETURNING id, article_id, user_id,created_at`, payload.Article.Id, payload.User.Id, time.Now()).Scan(
		&brk.Id, &brk.Article.Id, &brk.User.Id, &brk.CreatedAt,
	)

	if err != nil {
		return model.Bookmark{}, err
	}
	return brk, nil
}

func NewBookmarkRepository(database *sql.DB) BookmarkRepository {
	return &bookmarkRepository{db: database}
}
