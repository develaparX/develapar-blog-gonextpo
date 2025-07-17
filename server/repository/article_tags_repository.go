package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
)

type ArticleTagRepository interface {
	AssignTags(ctx context.Context, articleId int, tagId []int) error
	GetTagsByArticleId(ctx context.Context, articleId int) ([]model.Tags, error)
	GetArticleByTagId(ctx context.Context, tagId int) ([]model.Article, error)
	RemoveTagFromArticle(ctx context.Context, articleId, tagId int) error
}

type articleTagRepository struct {
	db *sql.DB
}

// RemoveTagFromArticle implements ArticleTagRepository.
func (a *articleTagRepository) RemoveTagFromArticle(ctx context.Context, articleId int, tagId int) error {
	_, err := a.db.ExecContext(ctx, `DELETE FROM article_tags WHERE article_id= $1 AND tag_id = $2`, articleId, tagId)

	return err
}

// AssignTags implements ArticleTagRepository.
func (a *articleTagRepository) AssignTags(ctx context.Context, articleId int, tagId []int) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM article_tags WHERE article_id = $1`, articleId)
	if err != nil {
		return err
	}

	for _, tagID := range tagId {
		_, err := tx.ExecContext(ctx, `INSERT INTO article_tags (article_id, tag_id) VALUES ($1,$2)`, articleId, tagID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// GetArticleByTagId implements ArticleTagRepository.
func (a *articleTagRepository) GetArticleByTagId(ctx context.Context, tagId int) ([]model.Article, error) {
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

	rows, err := a.db.QueryContext(ctx, query, tagId)
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
func (a *articleTagRepository) GetTagsByArticleId(ctx context.Context, articleId int) ([]model.Tags, error) {
	query := `SELECT t.id, t.name FROM tags t JOIN article_tags at ON t.id = at.tag_id WHERE at.article_id = $1`

	rows, err := a.db.QueryContext(ctx, query, articleId)
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
