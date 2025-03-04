package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ICourseContentRepository interface {
	CreateVideo(ctx context.Context, video *entity.CourseVideo) error
	UpdateVideo(ctx context.Context, id uuid.UUID, updates dto.CourseVideoUpdate) error
	DeleteVideo(ctx context.Context, id uuid.UUID) error
	GetVideoByID(ctx context.Context, id uuid.UUID) (*entity.CourseVideo, error)

	CreateMaterial(ctx context.Context, material *entity.CourseMaterial) error
	UpdateMaterial(ctx context.Context, id uuid.UUID, updates dto.CourseMaterialUpdate) error
	DeleteMaterial(ctx context.Context, id uuid.UUID) error
	GetMaterialByID(ctx context.Context, id uuid.UUID) (*entity.CourseMaterial, error)

	GetCourseContents(ctx context.Context, courseID uuid.UUID) ([]*entity.CourseVideo, []*entity.CourseMaterial, error)
}

type ICourseContentService interface {
	CreateVideo(ctx context.Context, courseID uuid.UUID,
		req dto.CreateCourseVideoRequest) (dto.CreateCourseVideoResponse, error)
	UpdateVideo(ctx context.Context, id uuid.UUID, req dto.UpdateCourseVideoRequest) error
	DeleteVideo(ctx context.Context, id uuid.UUID) error
	GetVideoUploadURLs(ctx context.Context, id uuid.UUID) (string, string, error) // videoURL, thumbnailURL, error

	CreateMaterial(ctx context.Context, courseID uuid.UUID,
		req dto.CreateCourseMaterialRequest) (dto.CreateCourseMaterialResponse, error)
	UpdateCourseMaterial(ctx context.Context, id uuid.UUID, req dto.UpdateCourseMaterialRequest) error
	DeleteCourseMaterial(ctx context.Context, id uuid.UUID) error
	GetMaterialUploadURL(ctx context.Context, id uuid.UUID) (string, error)

	GetCourseContents(ctx context.Context, courseID uuid.UUID) ([]*dto.CourseContentResponse, error)
}
