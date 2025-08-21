package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"develapar-server/model/dto"
	"time"

	"github.com/google/uuid"
)

type ArticleRepository interface {
	GetAll(ctx context.Context) ([]dto.ArticleResponse, error)
	GetAllWithPagination(ctx context.Context, offset, limit int) ([]model.Article, int, error)
	CreateArticle(ctx context.Context, payload model.Article) (model.Article, error)
	UpdateArticle(ctx context.Context, article model.Article) (model.Article, error)
	GetArticleById(ctx context.Context, id uuid.UUID) (model.Article, error)
	GetArticleByUserId(ctx context.Context, userId uuid.UUID) ([]model.Article, error)
	GetArticleByUserIdWithPagination(ctx context.Context, userId uuid.UUID, offset, limit int) ([]model.Article, int, error)
	GetArticleBySlug(ctx context.Context, slug string) (model.Article, error)
	GetArticleByCategory(ctx context.Context, cat string) ([]model.Article, error)
	GetArticleByCategoryWithPagination(ctx context.Context, cat string, offset, limit int) ([]model.Article, int, error)
	DeleteArticle(ctx context.Context, id uuid.UUID) error
}

type articleRepository struct {
	db *sql.DB
}

// GetArticleByCategory implements ArticleRepository.
func (a *articleRepository) GetArticleByCategory(ctx context.Context, cat string) ([]model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE c.name = $1;
	`
	rows, err := a.db.QueryContext(ctx, query, cat)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article

	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, err
		}

		article.User = &user
		article.Category = &category
		articles = append(articles, article)
	}

	return articles, nil
}

// GetArticleBySlug implements ArticleRepository.
func (a *articleRepository) GetArticleBySlug(ctx context.Context, slug string) (model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE a.slug = $1;
	`

	row := a.db.QueryRowContext(ctx, query, slug)

	var article model.Article
	var user model.User
	var category model.Category

	err := row.Scan(
		&article.Id, &article.Title, &article.Slug, &article.Content,
		&article.UserId, &article.CategoryId, &article.Views, &article.Status,
		&article.CreatedAt, &article.UpdatedAt,
		&user.Id, &user.Name, &user.Email, &user.Role,
		&category.Id, &category.Name,
	)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, err
	}

	article.User = &user
	article.Category = &category
	return article, nil
}

// DeleteArticle implements ArticleRepository.
func (a *articleRepository) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	_, err := a.db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return err
	}

	return nil
}

// GetArticleById implements ArticleRepository.
func (a *articleRepository) GetArticleById(ctx context.Context, id uuid.UUID) (model.Article, error) {
	query := `
		SELECT id, title, slug, content, views, user_id, category_id, status, created_at, updated_at
		FROM articles
		WHERE id = $1
	`
	var arc model.Article
	err := a.db.QueryRowContext(ctx, query, id).Scan(
		&arc.Id, &arc.Title, &arc.Slug, &arc.Content, &arc.Views,
		&arc.UserId, &arc.CategoryId, &arc.Status, &arc.CreatedAt, &arc.UpdatedAt,
	)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, err
	}
	return arc, nil
}

// GetArticleByUserId implements ArticleRepository.
func (a *articleRepository) GetArticleByUserId(ctx context.Context, userId uuid.UUID) ([]model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE a.user_id = $1
	ORDER BY a.created_at DESC;
	`

	rows, err := a.db.QueryContext(ctx, query, userId)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article

	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, err
		}

		article.User = &user
		article.Category = &category
		articles = append(articles, article)
	}

	return articles, nil
}

// UpdateArticle implements ArticleRepository.
func (a *articleRepository) UpdateArticle(ctx context.Context, article model.Article) (model.Article, error) {
	query := `
	UPDATE articles
	SET title = $1, slug = $2, content = $3, category_id = $4, status = $5, updated_at = NOW()
	WHERE id = $6
 	RETURNING id, title, slug, content, user_id, category_id, views, status, created_at, updated_at
	`
	row := a.db.QueryRowContext(ctx, query, article.Title, article.Slug, article.Content, article.CategoryId, article.Status, article.Id)
	var updated model.Article
	err := row.Scan(
		&updated.Id,
		&updated.Title,
		&updated.Slug,
		&updated.Content,
		&updated.UserId,
		&updated.CategoryId,
		&updated.Views,
		&updated.Status,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, err
	}

	return updated, nil
}

// GetAll implements ArticleRepository.
func (a *articleRepository) GetAll(ctx context.Context) ([]dto.ArticleResponse, error) {
	query := `
    SELECT 
        a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
        u.id, u.name, u.email, u.role,
        c.id, c.name,
        t.id, t.name
    FROM articles a
    JOIN users u ON a.user_id = u.id
    JOIN categories c ON a.category_id = c.id
    LEFT JOIN article_tags at ON a.id = at.article_id
    LEFT JOIN tags t ON at.tag_id = t.id
    ORDER BY a.created_at DESC, a.id;
    `

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articleMap := make(map[uuid.UUID]*dto.ArticleResponse)
	var articles []dto.ArticleResponse

	for rows.Next() {
		var article dto.ArticleResponse
		var user model.User
		var category model.Category

		var tagID sql.NullString
		var tagName sql.NullString

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
			&tagID, &tagName, // Scan kolom tag
		)
		if err != nil {
			return nil, err
		}

		// Cek apakah artikel ini sudah ada di map menggunakan ID UUID-nya
		if _, ok := articleMap[article.Id]; !ok {
			// Jika belum ada, ini adalah pertama kalinya kita melihat artikel ini
			article.User = &user
			article.Category = &category
			article.Tags = []model.Tags{} // Inisialisasi slice tags

			articles = append(articles, article)
			articleMap[article.Id] = &articles[len(articles)-1] // Simpan pointer-nya di map
		}

		// Tambahkan tag ke artikel yang sesuai jika tag-nya valid
		if tagID.Valid && tagName.Valid {
			// Parse UUID dari string yang didapat dari database
			parsedTagID, err := uuid.Parse(tagID.String)
			if err != nil {
				// Anda bisa mencatat error ini jika perlu, lalu lanjut ke baris berikutnya
				continue
			}

			tag := model.Tags{
				Id:   parsedTagID,
				Name: tagName.String,
			}

			// Ambil artikel dari map dan tambahkan tag baru
			existingArticle := articleMap[article.Id]
			existingArticle.Tags = append(existingArticle.Tags, tag)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

// CreateArticle implements ArticleRepository.
func (a *articleRepository) CreateArticle(ctx context.Context, payload model.Article) (model.Article, error) {
	var arc model.Article
	err := a.db.QueryRowContext(ctx, `
  INSERT INTO articles (id, title, content, slug, user_id, category_id, status, created_at, updated_at) 
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
  RETURNING id, title, slug, content, user_id, category_id, views, status, created_at, updated_at
`,
		payload.Id,
		payload.Title,
		payload.Content,
		payload.Slug,
		payload.User.Id,
		payload.Category.Id,
		payload.Status,
		time.Now(),
		time.Now(),
	).Scan(
		&arc.Id,
		&arc.Title,
		&arc.Slug,
		&arc.Content,
		&arc.UserId,
		&arc.CategoryId,
		&arc.Views,
		&arc.Status,
		&arc.CreatedAt,
		&arc.UpdatedAt,
	)

	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.Article{}, ctx.Err()
		}
		return model.Article{}, err
	}
	return arc, nil
}

// GetAllWithPagination implements ArticleRepository.
func (a *articleRepository) GetAllWithPagination(ctx context.Context, offset, limit int) ([]model.Article, int, error) {
	// First get the total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM articles`
	err := a.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}

	// Then get the paginated results
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	ORDER BY a.created_at DESC
	LIMIT $1 OFFSET $2;
	`

	rows, err := a.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		article.User = &user
		article.Category = &category
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return articles, totalCount, nil
}

// GetArticleByUserIdWithPagination implements ArticleRepository.
func (a *articleRepository) GetArticleByUserIdWithPagination(ctx context.Context, userId uuid.UUID, offset, limit int) ([]model.Article, int, error) {
	// First get the total count for this user
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM articles WHERE user_id = $1`
	err := a.db.QueryRowContext(ctx, countQuery, userId).Scan(&totalCount)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}

	// Then get the paginated results
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE a.user_id = $1
	ORDER BY a.created_at DESC
	LIMIT $2 OFFSET $3;
	`

	rows, err := a.db.QueryContext(ctx, query, userId, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		article.User = &user
		article.Category = &category
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return articles, totalCount, nil
}

// GetArticleByCategoryWithPagination implements ArticleRepository.
func (a *articleRepository) GetArticleByCategoryWithPagination(ctx context.Context, cat string, offset, limit int) ([]model.Article, int, error) {
	// First get the total count for this category
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM articles a JOIN categories c ON a.category_id = c.id WHERE c.name = $1`
	err := a.db.QueryRowContext(ctx, countQuery, cat).Scan(&totalCount)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}

	// Then get the paginated results
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.user_id, a.category_id, a.views, a.status, a.created_at, a.updated_at,
		u.id, u.name, u.email, u.role,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE c.name = $1
	ORDER BY a.created_at DESC
	LIMIT $2 OFFSET $3;
	`

	rows, err := a.db.QueryContext(ctx, query, cat, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.UserId, &article.CategoryId, &article.Views, &article.Status,
			&article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email, &user.Role,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		article.User = &user
		article.Category = &category
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return articles, totalCount, nil
}

func NewArticleRepository(database *sql.DB) ArticleRepository {
	return &articleRepository{db: database}
}
