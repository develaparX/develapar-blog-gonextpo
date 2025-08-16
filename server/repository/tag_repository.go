package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"time"

	"github.com/google/uuid"
)

type TagRepository interface {
	CreateTag(ctx context.Context, payload model.Tags) (model.Tags, error)
	GetAllTag(ctx context.Context) ([]model.Tags, error)
	GetTagById(ctx context.Context, id uuid.UUID) (model.Tags, error)
	GetTagByName(ctx context.Context, name string) (model.Tags, error)
	UpdateTag(ctx context.Context, payload model.Tags) (model.Tags, error)
	DeleteTag(ctx context.Context, id uuid.UUID) error
}

type tagRepository struct {
	db *sql.DB
}

// GetTagByName implements TagRepository.
func (t *tagRepository) GetTagByName(ctx context.Context, name string) (model.Tags, error) {
	var tag model.Tags

	err := t.db.QueryRowContext(ctx, `SELECT id, name, created_at, updated_at FROM tags WHERE name = $1`, name).Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return model.Tags{}, err
	}

	return tag, nil
}

// CreateTag implements TagRepository.
func (t *tagRepository) CreateTag(ctx context.Context, payload model.Tags) (model.Tags, error) {
	newId := uuid.Must(uuid.NewV7())
	var tag model.Tags
	err := t.db.QueryRowContext(ctx, `INSERT INTO tags (id, name, created_at, updated_at) VALUES($1, $2, $3, $4) RETURNING id, name, created_at, updated_at`, newId, payload.Name, time.Now(), time.Now()).Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return model.Tags{}, err
	}
	return tag, nil
}

// GetAllTag implements TagRepository.
func (t *tagRepository) GetAllTag(ctx context.Context) ([]model.Tags, error) {
	var listTag []model.Tags

	rows, err := t.db.QueryContext(ctx, `SELECT id, name, created_at, updated_at FROM tags`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag model.Tags

		err := rows.Scan(
			&tag.Id,
			&tag.Name,
			&tag.CreatedAt,
			&tag.UpdatedAt,
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
func (t *tagRepository) GetTagById(ctx context.Context, id uuid.UUID) (model.Tags, error) {
	var tag model.Tags

	err := t.db.QueryRowContext(ctx, `SELECT id, name, created_at, updated_at FROM tags WHERE id = $1`, id).Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return model.Tags{}, err
	}

	return tag, nil
}

// UpdateTag implements TagRepository.
func (t *tagRepository) UpdateTag(ctx context.Context, payload model.Tags) (model.Tags, error) {
	var tag model.Tags
	err := t.db.QueryRowContext(ctx, `UPDATE tags SET name = $1, updated_at = $2 WHERE id = $3 RETURNING id, name, created_at, updated_at`, payload.Name, time.Now(), payload.Id).Scan(&tag.Id, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return model.Tags{}, err
	}
	return tag, nil
}

// DeleteTag implements TagRepository.
func (t *tagRepository) DeleteTag(ctx context.Context, id uuid.UUID) error {
	_, err := t.db.ExecContext(ctx, `DELETE FROM tags WHERE id = $1`, id)
	return err
}

func NewTagRepository(database *sql.DB) TagRepository {
	return &tagRepository{db: database}
}
