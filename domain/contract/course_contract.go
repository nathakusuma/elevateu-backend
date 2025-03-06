package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type ICourseRepository interface {
	CreateCourse(ctx context.Context, course *entity.Course) error
	GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error)
	GetCourses(ctx context.Context, query dto.GetCoursesQuery,
		pageReq dto.PaginationRequest) ([]*entity.Course, dto.PaginationResponse, error)
	UpdateCourse(ctx context.Context, txWrapper database.ITransaction, updates *dto.CourseUpdate) error
	DeleteCourse(ctx context.Context, txWrapper database.ITransaction, id uuid.UUID) error

	CreateEnrollment(ctx context.Context, courseID, studentID uuid.UUID) error
	GetEnrolledCourses(ctx context.Context, studentID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*entity.Course, dto.PaginationResponse, error)
	GetEnrollment(ctx context.Context, courseID, studentID uuid.UUID) (*entity.CourseEnrollment, error)
}

type ICourseService interface {
	CreateCourse(ctx context.Context, req *dto.CreateCourseRequest) (dto.CreateCourseResponse, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error)
	GetCourses(ctx context.Context, query dto.GetCoursesQuery,
		paginationReq dto.PaginationRequest) ([]*dto.CourseResponse, dto.PaginationResponse, error)
	UpdateCourse(ctx context.Context, id uuid.UUID, req *dto.UpdateCourseRequest) error
	DeleteCourse(ctx context.Context, id uuid.UUID) error

	GetPreviewVideoUploadURL(ctx context.Context, id uuid.UUID) (string, error)

	CreateEnrollment(ctx context.Context, courseID, studentID uuid.UUID) error
	GetEnrolledCourses(ctx context.Context, studentID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*dto.CourseResponse, dto.PaginationResponse, error)
}
