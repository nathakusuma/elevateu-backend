package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type courseService struct {
	repo     contract.ICourseRepository
	fileUtil fileutil.IFileUtil
	uuid     uuidpkg.IUUID
}

func NewCourseService(
	repo contract.ICourseRepository,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.ICourseService {
	return &courseService{
		repo:     repo,
		fileUtil: fileUtil,
		uuid:     uuid,
	}
}

func (s *courseService) CreateCourse(ctx context.Context,
	req *dto.CreateCourseRequest) (dto.CreateCourseResponse, error) {
	courseID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[CourseService][CreateCourse] Failed to generate course ID")
		return dto.CreateCourseResponse{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	teacherAvatarURL, err := s.fileUtil.ValidateAndUploadFile(ctx, req.TeacherAvatar, fileutil.ImageContentTypes,
		fmt.Sprintf("courses/teacher_avatar/%s", courseID))
	if err != nil {
		return dto.CreateCourseResponse{}, err
	}

	thumbnailURL, err := s.fileUtil.ValidateAndUploadFile(ctx, req.Thumbnail, fileutil.ImageContentTypes,
		fmt.Sprintf("courses/thumbnail/%s", courseID))
	if err != nil {
		return dto.CreateCourseResponse{}, err
	}

	previewVideoPath := fmt.Sprintf("courses/preview_video/%s", courseID)
	previewVideoUploadURL, err := s.fileUtil.GetUploadSignedURL(previewVideoPath)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"path":  previewVideoPath,
		}, "[CourseService][CreateCourse] Failed to get preview video signed upload URL")
		return dto.CreateCourseResponse{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	course := &entity.Course{
		ID:               courseID,
		CategoryID:       req.CategoryID,
		Title:            req.Title,
		Description:      req.Description,
		TeacherName:      req.TeacherName,
		TeacherAvatarURL: teacherAvatarURL,
		ThumbnailURL:     thumbnailURL,
		PreviewVideoURL:  s.fileUtil.GetFullURL(previewVideoPath),
	}

	// Create course
	err = s.repo.CreateCourse(ctx, course)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[CourseService][CreateCourse] Failed to create course")
		return dto.CreateCourseResponse{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return dto.CreateCourseResponse{
		Course:                &dto.CourseResponse{ID: courseID},
		PreviewVideoUploadURL: previewVideoUploadURL,
	}, nil
}

func (s *courseService) GetCourseByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "course not found") {
			return nil, errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CourseService][GetCourseByID] Failed to get course by ID")
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	resp := &dto.CourseResponse{}
	err = resp.PopulateFromEntity(course, s.fileUtil.GetSignedURL)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":  err,
			"course": course,
		}, "[CourseService][GetCourseByID] Failed to populate course response from entity")
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return resp, nil
}

func (s *courseService) GetCourses(ctx context.Context, query dto.GetCoursesQuery,
	pageReq dto.PaginationRequest) ([]*dto.CourseResponse, dto.PaginationResponse, error) {
	courses, pageResp, err := s.repo.GetCourses(ctx, query, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"query": query,
			"page":  pageReq,
		}, "[CourseService][GetCourses] Failed to get courses")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	resp := make([]*dto.CourseResponse, len(courses))
	for i, course := range courses {
		resp[i] = &dto.CourseResponse{}
		err = resp[i].PopulateFromEntity(course, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":  err,
				"course": course,
			}, "[CourseService][GetCourses] Failed to populate course response from entity")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer.WithTraceID(traceID)
		}
	}

	return resp, pageResp, nil
}

func (s *courseService) UpdateCourse(ctx context.Context, id uuid.UUID, req *dto.UpdateCourseRequest) error {
	updates := &dto.CourseUpdate{
		ID:          id,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		TeacherName: req.TeacherName,
	}

	tx, err := s.repo.BeginTx()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseService][UpdateCourse] Failed to begin transaction")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			log.Error(map[string]interface{}{
				"error": err,
			}, "[CourseService][UpdateCourse] Failed to rollback transaction")
		}
	}()

	err = s.repo.UpdateCourse(ctx, tx, updates)
	if err != nil {
		if err.Error() == "course not found" {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[CourseService][UpdateCourse] Failed to update course")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	if req.TeacherAvatar != nil {
		teacherAvatarURL, err := s.fileUtil.ValidateAndUploadFile(ctx, req.TeacherAvatar, fileutil.ImageContentTypes,
			fmt.Sprintf("courses/teacher_avatar/%s", id))
		if err != nil {
			return err
		}
		updates.TeacherAvatarURL = &teacherAvatarURL
	}

	if req.Thumbnail != nil {
		thumbnailURL, err := s.fileUtil.ValidateAndUploadFile(ctx, req.Thumbnail, fileutil.ImageContentTypes,
			fmt.Sprintf("courses/thumbnail/%s", id))
		if err != nil {
			return err
		}
		updates.ThumbnailURL = &thumbnailURL
	}

	err = tx.Commit()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseService][UpdateCourse] Failed to commit transaction")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	return nil
}

func (s *courseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	tx, err := s.repo.BeginTx()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseService][DeleteCourse] Failed to begin transaction")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	defer func() {
		err = tx.Rollback()
		if err != nil {
			log.Error(map[string]interface{}{
				"error": err,
			}, "[CourseService][DeleteCourse] Failed to rollback transaction")
		}
	}()

	err = s.repo.DeleteCourse(ctx, tx, id)
	if err != nil {
		if err.Error() == "course not found" {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CourseService][DeleteCourse] Failed to delete course")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("courses/teacher_avatar/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CourseService][DeleteCourse] Failed to delete teacher avatar")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("courses/thumbnail/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CourseService][DeleteCourse] Failed to delete thumbnail")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("courses/preview_video/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CourseService][DeleteCourse] Failed to delete preview video")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	err = tx.Commit()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseService][DeleteCourse] Failed to commit transaction")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	return nil
}

func (s *courseService) GetPreviewVideoUploadURL(_ context.Context, id uuid.UUID) (string, error) {
	url, err := s.fileUtil.GetUploadSignedURL(fmt.Sprintf("courses/preview_video/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[CourseService][GetPreviewVideoUploadURL] Failed to get preview video signed upload URL")
		return "", errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return url, nil
}
