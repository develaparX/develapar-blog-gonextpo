package repository

import (
	"database/sql"
	"develapar-server/model"
)

type ArticleRepository interface {
	GetAll() ([]model.Article, error)
	CreateArticle(payload model.Article) (model.Article, error)
}

type articleRepository struct {
	db *sql.DB
}

// GetAll implements ArticleRepository.
func (a *articleRepository) GetAll() ([]model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	ORDER BY a.created_at DESC;
	`

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article
		var user model.User
		var category model.Category

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content,
			&article.Views, &article.CreatedAt, &article.UpdatedAt,
			&user.Id, &user.Name, &user.Email,
			&category.Id, &category.Name,
		)
		if err != nil {
			return nil, err
		}

		article.User = user
		article.Category = category
		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

// CreateArticle implements ArticleRepository.
func (a *articleRepository) CreateArticle(payload model.Article) (model.Article, error) {
	var arc model.Article
	err := a.db.QueryRow(`
  INSERT INTO articles (title, content, slug, user_id, category_id) 
  VALUES ($1, $2, $3, $4, $5) 
  RETURNING id, title, slug, content, user_id, category_id, views, created_at, updated_at
`,
		payload.Title,
		payload.Content,
		payload.Slug,
		payload.User.Id,
		payload.Category.Id,
	).Scan(
		&arc.Id,
		&arc.Title,
		&arc.Slug,
		&arc.Content,
		&arc.User.Id,     // hanya ambil ID-nya
		&arc.Category.Id, // hanya ambil ID-nya
		&arc.Views,
		&arc.CreatedAt,
		&arc.UpdatedAt,
	)

	if err != nil {
		return model.Article{}, err
	}
	return arc, nil
}

func NewArticleRepository(database *sql.DB) ArticleRepository {
	return &articleRepository{db: database}
}
