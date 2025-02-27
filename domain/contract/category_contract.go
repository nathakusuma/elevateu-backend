package contract

import (
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ICategoryRepository interface {
	CreateCategory(uuid2 uuid.UUID, name string) error
	GetAllCategories() ([]entity.Category, error)
	UpdateCategory(id uuid.UUID, name string) error
	DeleteCategory(id uuid.UUID) error
}

type ICategoryService interface {
	CreateCategory(name string) (entity.Category, error)
	GetAllCategories() ([]entity.Category, error)
	UpdateCategory(id uuid.UUID, name string) error
	DeleteCategory(id uuid.UUID) error
}
