package repository

import (
	"database/sql"
	"develapar-server/model"
)

type TagRepository interface {
	CreateTag(payload model.Tags) (model.Tags, error)
	GetAllTag() ([]model.Tags, error)
	GetTagById(id int) (model.Tags, error)
	GetTagByName(name string) (model.Tags, error)
}

type tagRepository struct {
	db *sql.DB
}

// GetTagByName implements TagRepository.
func (t *tagRepository) GetTagByName(name string) (model.Tags, error) {
	var tag model.Tags

	err := t.db.QueryRow(`SELECT id, name FROM tags WHERE name = $1`, name).Scan(&tag.Id, &tag.Name)
	if err != nil {
		return model.Tags{}, err
	}

	return tag, nil
}

// CreateTag implements TagRepository.
func (t *tagRepository) CreateTag(payload model.Tags) (model.Tags, error) {
	var tag model.Tags
	err := t.db.QueryRow(`INSERT INTO tags (name) VALUES($1) RETURNING id, name`, payload.Name).Scan(&tag.Id, &tag.Name)
	if err != nil {
		return model.Tags{}, err
	}
	return tag, nil
}

// GetAllTag implements TagRepository.
func (t *tagRepository) GetAllTag() ([]model.Tags, error) {
	var listTag []model.Tags

	rows, err := t.db.Query(`SELECT id, name FROM tags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag model.Tags

		err := rows.Scan(
			&tag.Id,
			&tag.Name,
		)

		if err != nil {
			return nil, err
		}

		listTag = append(listTag, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listTag, nil
}

// GetTagById implements TagRepository.
func (t *tagRepository) GetTagById(id int) (model.Tags, error) {
	var tag model.Tags

	err := t.db.QueryRow(`SELECT id, name FROM tags WHERE id = $1`, id).Scan(&tag.Id, &tag.Name)
	if err != nil {
		return model.Tags{}, err
	}

	return tag, nil
}

func NewTagRepository(database *sql.DB) TagRepository {
	return &tagRepository{db: database}
}
