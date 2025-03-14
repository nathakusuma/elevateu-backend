package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type categoryService struct {
	repo contract.ICategoryRepository
	uuid uuidpkg.IUUID
}

func NewCategoryService(repo contract.ICategoryRepository, uuid uuidpkg.IUUID) contract.ICategoryService {
	return &categoryService{
		repo: repo,
		uuid: uuid,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, name string) (entity.Category, error) {
	id, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"name":  name,
		}, "Failed to generate category ID")
		return entity.Category{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err = s.repo.CreateCategory(ctx, id, name); err != nil {
		if strings.HasPrefix(err.Error(), "conflict name") {
			return entity.Category{}, errorpkg.ErrCategoryNameExists()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"name":  name,
		}, "Failed to create category")
		return entity.Category{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	category := entity.Category{ID: id}

	log.Info(ctx, map[string]interface{}{
		"category": category,
	}, "Category created")

	return category, nil
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repo.GetAllCategories(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to get all categories")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return categories, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id uuid.UUID, name string) error {
	if err := s.repo.UpdateCategory(ctx, id, name); err != nil {
		if strings.HasPrefix(err.Error(), "category not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
			"name":  name,
		}, "Failed to update category")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"id":   id,
		"name": name,
	}, "Category updated")

	return nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		if strings.HasPrefix(err.Error(), "category not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to delete category")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"id": id,
	}, "Category deleted")

	return nil
}
