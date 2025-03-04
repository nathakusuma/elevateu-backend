package dto

import (
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type CourseResponse struct {
	ID               uuid.UUID `json:"id,omitempty"`
	Category         string    `json:"category,omitempty"`
	Title            string    `json:"title,omitempty"`
	Description      string    `json:"description,omitempty"`
	TeacherName      string    `json:"teacher_name,omitempty"`
	TeacherAvatarURL string    `json:"teacher_avatar_url,omitempty"`
	ThumbnailURL     string    `json:"thumbnail_url,omitempty"`
	PreviewVideoURL  string    `json:"preview_video_url,omitempty"`
	Rating           *float64  `json:"rating,omitempty"`
	RatingCount      *int64    `json:"rating_count,omitempty"`
	EnrollmentCount  *int64    `json:"enrollment_count,omitempty"`
	ContentCount     *int      `json:"content_count,omitempty"`
	TotalDuration    *int      `json:"total_duration,omitempty"`
}

func (c *CourseResponse) PopulateFromEntity(course *entity.Course,
	urlSigner func(string) (string, error)) error {
	var err error
	c.TeacherAvatarURL, err = urlSigner(fmt.Sprintf("courses/teacher_avatar/%s", course.ID))
	if err != nil {
		return fmt.Errorf("failed to sign teacher avatar URL: %w", err)
	}

	c.ThumbnailURL, err = urlSigner(fmt.Sprintf("courses/thumbnail/%s", course.ID))
	if err != nil {
		return fmt.Errorf("failed to sign thumbnail URL: %w", err)
	}

	c.PreviewVideoURL, err = urlSigner(fmt.Sprintf("courses/preview_video/%s", course.ID))
	if err != nil {
		return fmt.Errorf("failed to sign preview video URL: %w", err)
	}

	if course.Category != nil {
		c.Category = course.Category.Name
	}

	c.ID = course.ID
	c.Title = course.Title
	c.Description = course.Description
	c.Rating = &course.Rating
	c.RatingCount = &course.RatingCount
	c.TeacherName = course.TeacherName
	c.EnrollmentCount = &course.EnrollmentCount
	c.ContentCount = &course.ContentCount
	c.TotalDuration = &course.TotalDuration

	return nil
}

type CreateCourseRequest struct {
	CategoryID    uuid.UUID `form:"category_id" validate:"required,uuid"`
	Title         string    `form:"title" validate:"required,min=3,max=50"`
	Description   string    `form:"description" validate:"required,min=3,max=1000"`
	TeacherName   string    `form:"teacher_name" validate:"required,min=3,max=50"`
	TeacherAvatar *multipart.FileHeader
	Thumbnail     *multipart.FileHeader
}

type CreateCourseResponse struct {
	Course                *CourseResponse `json:"course"`
	PreviewVideoUploadURL string          `json:"preview_video_upload_url"`
}

type GetCoursesQuery struct {
	CategoryID uuid.UUID `query:"category_id" validate:"omitempty,uuid"`
	Title      string    `query:"title" validate:"omitempty"`
}

type CourseUpdate struct {
	ID          uuid.UUID  `db:"id"`
	CategoryID  *uuid.UUID `db:"category_id"`
	Title       *string    `db:"title"`
	Description *string    `db:"description"`
	TeacherName *string    `db:"teacher_name"`
}

type UpdateCourseRequest struct {
	CategoryID    *uuid.UUID `form:"category_id" validate:"omitempty,uuid"`
	Title         *string    `form:"title" validate:"omitempty,min=3,max=50"`
	Description   *string    `form:"description" validate:"omitempty,min=3,max=1000"`
	TeacherName   *string    `form:"teacher_name" validate:"omitempty,min=3,max=60"`
	TeacherAvatar *multipart.FileHeader
	Thumbnail     *multipart.FileHeader
}
