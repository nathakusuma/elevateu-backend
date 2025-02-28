package contract

import (
	"context"
	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ICourseRepository interface {
	BeginTx() (*sqlx.Tx, error)
	CreateCourse(ctx context.Context, course *entity.Course) error
	GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetCourses(ctx context.Context, query dto.GetCoursesQuery,
		pageReq dto.PaginationRequest) ([]*entity.Course, dto.PaginationResponse, error)
	UpdateCourse(ctx context.Context, tx sqlx.ExtContext, updates *dto.CourseUpdate) error
	DeleteCourse(ctx context.Context, tx sqlx.ExtContext, id uuid.UUID) error
}

type ICourseService interface {
	CreateCourse(ctx context.Context, req *dto.CreateCourseRequest) (dto.CreateCourseResponse, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error)
	GetCourses(ctx context.Context, query dto.GetCoursesQuery,
		paginationReq dto.PaginationRequest) ([]*dto.CourseResponse, dto.PaginationResponse, error)
	UpdateCourse(ctx context.Context, id uuid.UUID, req *dto.UpdateCourseRequest) error
	DeleteCourse(ctx context.Context, id uuid.UUID) error

	GetPreviewVideoUploadURL(ctx context.Context, id uuid.UUID) (string, error)
}
