package repository

import (
	"database/sql"
	"develapar-server/model"
)

type CategoryRepository interface {
	GetAll() ([]model.Category, error)
	CreateCategory(payload model.Category) (model.Category, error)
}

type categoryRepository struct {
	db *sql.DB
}

// CreateCategory implements CategoryRepository.
func (c *categoryRepository) CreateCategory(payload model.Category) (model.Category, error) {
	var cat model.Category
	err := c.db.QueryRow(`INSERT INTO categories (name) VALUES($1) RETURNING id, name`, payload.Name).Scan(&cat.Id, &cat.Name)
	if err != nil {
		return model.Category{}, err
	}
	return cat, nil
}

// GetAll implements CategoryRepository.
func (c *categoryRepository) GetAll() ([]model.Category, error) {
	var listCategory []model.Category

	rows, err := c.db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.Id,
			&category.Name,
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
