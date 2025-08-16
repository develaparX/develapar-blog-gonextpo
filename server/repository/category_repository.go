package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"time"

	"github.com/google/uuid"
)

type CategoryRepository interface {
	GetAll(ctx context.Context) ([]model.Category, error)
	CreateCategory(ctx context.Context, payload model.Category) (model.Category, error)
	GetCategoryById(ctx context.Context, id uuid.UUID) (model.Category, error)
	UpdateCategory(ctx context.Context, payload model.Category) (model.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type categoryRepository struct {
	db *sql.DB
}

// GetCategoryById implements CategoryRepository.
func (c *categoryRepository) GetCategoryById(ctx context.Context, id uuid.UUID) (model.Category, error) {
	query := `
	SELECT id, name, created_at, updated_at
	FROM categories
	WHERE id = $1
	`

	var cat model.Category
	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&cat.Id, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt,
	)
	if err != nil {
		return model.Category{}, err
	}

	return cat, nil
}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepository) UpdateCategory(ctx context.Context, payload model.Category) (model.Category, error) {
	var cat model.Category
	err := c.db.QueryRowContext(ctx, `UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3 RETURNING id, name, created_at, updated_at`, payload.Name, time.Now(), payload.Id).Scan(&cat.Id, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)

	if err != nil {
		return model.Category{}, err
	}

	return cat, nil
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(ctx context.Context, payload model.Category) (model.Category, error) {
	newId := uuid.Must(uuid.NewV7())
	var cat model.Category
	err := c.db.QueryRowContext(ctx, `INSERT INTO categories (id, name, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id, name, created_at, updated_at`, newId, payload.Name, time.Now(), time.Now()).Scan(&cat.Id, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt)
	if err != nil {
		return model.Category{}, err
	}
	return cat, nil
}

// GetAll implements CategoryRepository.
func (c *categoryRepository) GetAll(ctx context.Context) ([]model.Category, error) {
	var listCategory []model.Category

	rows, err := c.db.QueryContext(ctx, `SELECT id, name, created_at, updated_at FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.Id,
			&category.Name,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		listCategory = append(listCategory, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listCategory, nil
}

func NewCategoryRepository(database *sql.DB) CategoryRepository {
	return &categoryRepository{db: database}
}
