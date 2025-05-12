package repository

import (
	"database/sql"
	"develapar-server/model"
)

type ArticleRepository interface {
	GetAll() ([]model.Article, error)
	CreateArticle(payload model.Article) (model.Article, error)
	UpdateArticle(article model.Article) (model.Article, error)
	GetArticleById(id int) (model.Article, error)
	GetArticleByUserId(userId int) ([]model.Article, error)
	GetArticleBySlug(slug string) (model.Article, error)
	GetArticleByCategory(cat string) ([]model.Article, error)
	DeleteArticle(id int) error
}

type articleRepository struct {
	db *sql.DB
}

// GetArticleByCategory implements ArticleRepository.
func (a *articleRepository) GetArticleByCategory(cat string) ([]model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE c.name = $1;
	`
	rows, err := a.db.Query(query, cat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article

	for rows.Next() {
		var a model.Article
		var u model.User
		var c model.Category

		err := rows.Scan(
			&a.Id, &a.Title, &a.Slug, &a.Content, &a.Views, &a.CreatedAt, &a.UpdatedAt,
			&u.Id, &u.Name, &u.Email,
			&c.Id, &c.Name,
		)
		if err != nil {
			return nil, err
		}

		a.User = u
		a.Category = c
		articles = append(articles, a)
	}

	return articles, nil

}

// GetArticleBySlug implements ArticleRepository.
func (a *articleRepository) GetArticleBySlug(slug string) (model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE a.slug = $1;
	`

	row := a.db.QueryRow(query, slug)

	var article model.Article
	var user model.User
	var category model.Category

	err := row.Scan(
		&article.Id, &article.Title, &article.Slug, &article.Content, &article.Views, &article.CreatedAt, &article.UpdatedAt,
		&user.Id, &user.Name, &user.Email,
		&category.Id, &category.Name,
	)
	if err != nil {
		return model.Article{}, err
	}

	article.User = user
	article.Category = category
	return article, nil
}

// DeleteArticle implements ArticleRepository.
func (a *articleRepository) DeleteArticle(id int) error {
	_, err := a.db.Exec(`DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

// GetArticleById implements ArticleRepository.
func (a *articleRepository) GetArticleById(id int) (model.Article, error) {
	query := `
		SELECT id, title, slug, content, views, user_id, category_id, created_at, updated_at
		FROM articles
		WHERE id = $1
	`
	var arc model.Article
	err := a.db.QueryRow(query, id).Scan(
		&arc.Id, &arc.Title, &arc.Slug, &arc.Content, &arc.Views,
		&arc.User.Id, &arc.Category.Id, &arc.CreatedAt, &arc.UpdatedAt,
	)
	return arc, err
}

// GetArticleByUserId implements ArticleRepository.
func (a *articleRepository) GetArticleByUserId(userId int) ([]model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	WHERE a.user_id = $1
	ORDER BY a.created_at DESC;
	`

	rows, err := a.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article

	for rows.Next() {
		var a model.Article
		var u model.User
		var c model.Category

		err := rows.Scan(
			&a.Id, &a.Title, &a.Slug, &a.Content, &a.Views, &a.CreatedAt, &a.UpdatedAt,
			&u.Id, &u.Name, &u.Email,
			&c.Id, &c.Name,
		)
		if err != nil {
			return nil, err
		}

		a.User = u
		a.Category = c
		articles = append(articles, a)
	}

	return articles, nil
}

// UpdateArticle implements ArticleRepository.
func (a *articleRepository) UpdateArticle(article model.Article) (model.Article, error) {
	query := `
	UPDATE articles
	SET title = $1, slug = $2, content = $3, category_id=$4, updated_at = NOW()
	WHERE id = $5
 	RETURNING id, title, slug, content, category_id, created_at, updated_at
	`
	row := a.db.QueryRow(query, article.Title, article.Slug, article.Content, article.Category.Id, article.Id)
	var updated model.Article
	err := row.Scan(
		&updated.Id,
		&updated.Title,
		&updated.Slug,
		&updated.Content,
		&updated.Category.Id,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return model.Article{}, err
	}

	return updated, nil
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
