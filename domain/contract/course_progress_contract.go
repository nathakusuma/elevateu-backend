package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type ICourseProgressRepository interface {
	UpdateVideoProgress(ctx context.Context, txWrapper database.ITransaction,
		progress entity.CourseVideoProgress) (bool, error)
	UpdateMaterialProgress(ctx context.Context, txWrapper database.ITransaction,
		progress entity.CourseMaterialProgress) (bool, error)
	IncrementCourseProgress(ctx context.Context, txWrapper database.ITransaction, courseID,
		studentID uuid.UUID) (bool, error)

	GetContentCourseID(ctx context.Context, contentID uuid.UUID, contentType string) (uuid.UUID, error)

	BatchDecrementCourseProgress(ctx context.Context, txWrapper database.ITransaction, courseID uuid.UUID,
		contentID uuid.UUID, contentType string) error
}

type ICourseProgressService interface {
	UpdateVideoProgress(ctx context.Context, studentID, videoID uuid.UUID,
		req dto.UpdateCourseVideoProgressRequest) error
	UpdateMaterialProgress(ctx context.Context, studentID uuid.UUID, materialID uuid.UUID) error
}
