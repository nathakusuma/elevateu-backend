package contract

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ICourseFeedbackRepository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)

	CreateFeedback(ctx context.Context, tx sqlx.ExtContext, feedback *entity.CourseFeedback) error
	GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*entity.CourseFeedback, dto.PaginationResponse, error)
	GetFeedbackByID(ctx context.Context, feedbackID uuid.UUID) (*entity.CourseFeedback, error)
	UpdateFeedback(ctx context.Context, tx sqlx.ExtContext, feedbackID uuid.UUID, updates dto.CourseFeedbackUpdate) error
	DeleteFeedback(ctx context.Context, tx sqlx.ExtContext, feedbackID uuid.UUID) error

	UpdateCourseRating(ctx context.Context, tx sqlx.ExtContext, courseID uuid.UUID, count int64, rating,
		total float64) error
}

type ICourseFeedbackService interface {
	CreateFeedback(ctx context.Context, courseID uuid.UUID, req dto.CreateCourseFeedbackRequest) error
	GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*dto.CourseFeedbackResponse, dto.PaginationResponse, error)
	UpdateFeedback(ctx context.Context, feedbackID uuid.UUID, req dto.UpdateCourseFeedbackRequest) error
	DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error
}
