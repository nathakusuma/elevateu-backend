package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ICourseProgressRepository interface {
	BeginTx() (*sqlx.Tx, error)

	UpdateVideoProgress(ctx context.Context, tx sqlx.ExtContext, progress entity.CourseVideoProgress) (bool, error)
	UpdateMaterialProgress(ctx context.Context, tx sqlx.ExtContext, progress entity.CourseMaterialProgress) (bool, error)
	IncrementCourseProgress(ctx context.Context, tx *sqlx.Tx, courseID, studentID uuid.UUID) (bool, error)

	GetContentCourseID(ctx context.Context, contentID uuid.UUID, contentType string) (uuid.UUID, error)

	BatchDecrementCourseProgress(ctx context.Context, tx sqlx.ExtContext, courseID uuid.UUID, contentID uuid.UUID,
		contentType string) error
}

type ICourseProgressService interface {
	UpdateVideoProgress(ctx context.Context, studentID, videoID uuid.UUID, req dto.UpdateCourseVideoProgressRequest) error
	UpdateMaterialProgress(ctx context.Context, studentID uuid.UUID, materialID uuid.UUID) error
}
