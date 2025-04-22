package repository

import (
	"database/sql"
	"develapar-server/model"
)

type ArticleTagRepository interface {
	AssignTags(articleId int, tagId []int) error
	GetTagsByArticleId(articleId int) ([]model.Tags, error)
	GetArticleByTagId(tagId int) ([]model.Article, error)
}

type articleTagRepository struct {
	db *sql.DB
}

// AssignTags implements ArticleTagRepository.
func (a *articleTagRepository) AssignTags(articleId int, tagId []int) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM article_tags WHERE article_id = $1`, articleId)
	if err != nil {
		return err
	}

	for _, tagID := range tagId {
		_, err := tx.Exec(`INSERT INTO article_tags (article_id, tag_id) VALUES ($1,$2)`, articleId, tagID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// GetArticleByTagId implements ArticleTagRepository.
func (a *articleTagRepository) GetArticleByTagId(tagId int) ([]model.Article, error) {
	query := `
	SELECT 
		a.id, a.title, a.slug, a.content, a.views, a.created_at, a.updated_at,
		u.id, u.name, u.email,
		c.id, c.name
	FROM articles a
	JOIN users u ON a.user_id = u.id
	JOIN categories c ON a.category_id = c.id
	JOIN article_tags at ON at.article_id = a.id
	WHERE at.tag_id = $1
	ORDER BY a.created_at DESC`

	rows, err := a.db.Query(query, tagId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []model.Article
	for rows.Next() {
		var article model.Article

		err := rows.Scan(
			&article.Id, &article.Title, &article.Slug, &article.Content, &article.Views, &article.CreatedAt, &article.UpdatedAt, &article.User.Id, &article.User.Name, &article.User.Email, &article.Category.Id, &article.Category.Name,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// GetTagsByArticleId implements ArticleTagRepository.
func (a *articleTagRepository) GetTagsByArticleId(articleId int) ([]model.Tags, error) {
	query := `SELECT t.id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = $1`

	rows, err := a.db.Query(query, articleId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tags

	for rows.Next() {
		var tag model.Tags
		err := rows.Scan(&tag.Id, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func NewArticleTagRepository(database *sql.DB) ArticleTagRepository {
	return &articleTagRepository{db: database}
}
