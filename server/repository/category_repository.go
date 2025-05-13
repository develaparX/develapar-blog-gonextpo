package repository

import (
	"database/sql"
	"develapar-server/model"
)

type CategoryRepository interface {
	GetAll() ([]model.Category, error)
	CreateCategory(payload model.Category) (model.Category, error)
	GetCategoryById(id int) (model.Category, error)
	UpdateCategory(payload model.Category) (model.Category, error)
	DeleteCategory(id int) error
}

type categoryRepository struct {
	db *sql.DB
}

// FindCategoryById implements CategoryRepository.
func (c *categoryRepository) GetCategoryById(id int) (model.Category, error) {
	query := `
	SELECT id, name
	FROM categories
	WHERE id = $1
	`

	var cat model.Category
	err := c.db.QueryRow(query, id).Scan(
		&cat.Id, &cat.Name,
	)
	if err != nil {
		return model.Category{}, err
	}

	return cat, nil

}

// DeleteCategory implements CategoryRepository.
func (c *categoryRepository) DeleteCategory(id int) error {
	_, err := c.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateCategory implements CategoryRepository.
func (c *categoryRepository) UpdateCategory(payload model.Category) (model.Category, error) {
	var cat model.Category
	err := c.db.QueryRow(`UPDATE categories SET name = $1 WHERE id = $2 RETURNING id, name`, payload.Name, payload.Id).Scan(&cat.Id, &cat.Name)

	if err != nil {
		return model.Category{}, err
	}

	return cat, nil
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
