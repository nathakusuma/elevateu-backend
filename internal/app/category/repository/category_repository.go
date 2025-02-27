package repository

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(conn *sqlx.DB) contract.ICategoryRepository {
	return &categoryRepository{
		db: conn,
	}
}

func (r *categoryRepository) CreateCategory(id uuid.UUID, name string) error {
	_, err := r.db.Exec(`INSERT INTO categories (id, name) VALUES ($1, $2)`, id, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "categories_name_key" {
			return fmt.Errorf("conflict name: %w", err)
		}

		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

func (r *categoryRepository) GetAllCategories() ([]entity.Category, error) {
	var categories []entity.Category
	err := r.db.Select(&categories, `SELECT * FROM categories`)
	return categories, err
}

func (r *categoryRepository) UpdateCategory(id uuid.UUID, name string) error {
	res, err := r.db.Exec(`UPDATE categories SET name = $1 WHERE id = $2`, name, id)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *categoryRepository) DeleteCategory(id uuid.UUID) error {
	res, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
