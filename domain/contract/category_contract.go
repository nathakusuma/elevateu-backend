package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ICategoryRepository interface {
	CreateCategory(ctx context.Context, id uuid.UUID, name string) error
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name string) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type ICategoryService interface {
	CreateCategory(ctx context.Context, name string) (entity.Category, error)
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
	UpdateCategory(ctx context.Context, id uuid.UUID, name string) error
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}
