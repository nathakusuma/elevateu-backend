package dto

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type CourseContentResponse struct {
	Type         string    `json:"type"`
	ID           uuid.UUID `json:"id"`
	URL          string    `json:"url,omitempty"`
	ThumbnailURL string    `json:"thumbnail_url,omitempty"`
	Title        string    `json:"title,omitempty"`
	Description  string    `json:"description,omitempty"`
	Subtitle     string    `json:"subtitle,omitempty"`
	Duration     int       `json:"duration,omitempty"`
	IsFree       bool      `json:"is_free,omitempty"`
}

func (c *CourseContentResponse) PopulateFromCourseVideo(video *entity.CourseVideo,
	urlSigner func(string) (string, error)) error {
	c.Type = "video"
	c.ID = video.ID
	c.Title = video.Title
	c.Description = video.Description
	c.Duration = video.Duration
	c.IsFree = video.IsFree

	var err error
	c.URL, err = urlSigner(fmt.Sprintf("course_videos/video/%s", video.ID.String()))
	if err != nil {
		return err
	}

	c.ThumbnailURL, err = urlSigner(fmt.Sprintf("course_videos/thumbnail/%s", video.ID.String()))
	if err != nil {
		return err
	}

	return nil
}

func (c *CourseContentResponse) PopulateFromCourseMaterial(material *entity.CourseMaterial,
	urlSigner func(string) (string, error)) error {
	c.Type = "material"
	c.ID = material.ID
	c.Title = material.Title
	c.Subtitle = material.Subtitle
	c.IsFree = material.IsFree

	var err error
	c.URL, err = urlSigner(fmt.Sprintf("course_materials/material/%s", material.ID.String()))
	if err != nil {
		return err
	}

	return nil
}

type CourseVideoUpdate struct {
	Title       *string `db:"title"`
	Description *string `db:"description"`
	Duration    *int    `db:"duration"`
	IsFree      *bool   `db:"is_free"`
	Order       *int    `db:"order"`
}

type CourseMaterialUpdate struct {
	Title    *string `db:"title"`
	Subtitle *string `db:"subtitle"`
	IsFree   *bool   `db:"is_free"`
	Order    *int    `db:"order"`
}

type CreateCourseVideoRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"required,min=3,max=1000"`
	Duration    int    `json:"duration" validate:"required,min=1"`
	IsFree      bool   `json:"is_free"`
	Order       int    `json:"order" validate:"required"`
}

type CreateCourseVideoResponse struct {
	CourseContent      *CourseContentResponse `json:"course_content,omitempty"`
	VideoUploadURL     string                 `json:"video_upload_url"`
	ThumbnailUploadURL string                 `json:"thumbnail_upload_url"`
}

type CreateCourseMaterialResponse struct {
	CourseContent     *CourseContentResponse `json:"course_content"`
	MaterialUploadURL string                 `json:"material_upload_url"`
}

type UpdateCourseVideoRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=50"`
	Description *string `json:"description" validate:"omitempty,min=3,max=1000"`
	Duration    *int    `json:"duration" validate:"omitempty,min=1"`
	IsFree      *bool   `json:"is_free"`
	Order       *int    `json:"order"`
}

type CreateCourseMaterialRequest struct {
	Title    string `json:"title" validate:"required,min=3,max=50"`
	Subtitle string `json:"subtitle" validate:"required,min=3,max=50"`
	IsFree   bool   `json:"is_free"`
	Order    int    `json:"order" validate:"required"`
}

type UpdateCourseMaterialRequest struct {
	Title    *string `form:"title" validate:"omitempty,min=3,max=50"`
	Subtitle *string `form:"subtitle" validate:"omitempty,min=3,max=50"`
	IsFree   *bool   `form:"is_free"`
	Order    *int    `form:"order"`
}
