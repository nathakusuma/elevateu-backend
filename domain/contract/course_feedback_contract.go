package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type ICourseFeedbackRepository interface {
	CreateFeedback(ctx context.Context, txWrapper database.ITransaction, feedback *entity.CourseFeedback) error
	GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*entity.CourseFeedback, dto.PaginationResponse, error)
	GetFeedbackByID(ctx context.Context, feedbackID uuid.UUID) (*entity.CourseFeedback, error)
	UpdateFeedback(ctx context.Context, txWrapper database.ITransaction, feedbackID uuid.UUID,
		updates dto.CourseFeedbackUpdate) error
	DeleteFeedback(ctx context.Context, txWrapper database.ITransaction, feedbackID uuid.UUID) error

	UpdateCourseRating(ctx context.Context, txWrapper database.ITransaction, courseID uuid.UUID, count int64, rating,
		total float64) error
}

type ICourseFeedbackService interface {
	CreateFeedback(ctx context.Context, courseID uuid.UUID, req dto.CreateCourseFeedbackRequest) error
	GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*dto.CourseFeedbackResponse, dto.PaginationResponse, error)
	UpdateFeedback(ctx context.Context, feedbackID uuid.UUID, req dto.UpdateCourseFeedbackRequest) error
	DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error
}
