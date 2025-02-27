package service

import (
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

func (s *categoryService) CreateCategory(name string) (entity.Category, error) {
	id, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"name":  name,
		}, "[CategoryService][CreateCategory] Failed to generate category ID")
		return entity.Category{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	if err = s.repo.CreateCategory(id, name); err != nil {
		if strings.HasPrefix("conflict name", err.Error()) {
			return entity.Category{}, errorpkg.ErrCategoryNameExists
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"name":  name,
		}, "[CategoryService][CreateCategory] Failed to create category")
		return entity.Category{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return entity.Category{ID: id}, nil
}

func (s *categoryService) GetAllCategories() ([]entity.Category, error) {
	categories, err := s.repo.GetAllCategories()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CategoryService][GetAllCategories] Failed to get all categories")
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return categories, nil
}

func (s *categoryService) UpdateCategory(id uuid.UUID, name string) error {
	if err := s.repo.UpdateCategory(id, name); err != nil {
		if strings.HasPrefix("category not found", err.Error()) {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
			"name":  name,
		}, "[CategoryService][UpdateCategory] Failed to update category")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return nil
}

func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	if err := s.repo.DeleteCategory(id); err != nil {
		if strings.HasPrefix("category not found", err.Error()) {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CategoryService][DeleteCategory] Failed to delete category")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return nil
}
