package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type courseService struct {
	repo      contract.ICourseRepository
	fileUtil  fileutil.IFileUtil
	txManager database.ITransactionManager
	uuid      uuidpkg.IUUID
}

func NewCourseService(
	repo contract.ICourseRepository,
	fileUtil fileutil.IFileUtil,
	txManager database.ITransactionManager,
	uuid uuidpkg.IUUID,
) contract.ICourseService {
	return &courseService{
		repo:      repo,
		fileUtil:  fileUtil,
		txManager: txManager,
		uuid:      uuid,
	}
}

func (s *courseService) CreateCourse(ctx context.Context,
	req *dto.CreateCourseRequest) (dto.CreateCourseResponse, error) {
	courseID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to generate course ID")
		return dto.CreateCourseResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	_, err = s.fileUtil.ValidateAndUploadFile(ctx, req.TeacherAvatar, fileutil.ImageContentTypes,
		fmt.Sprintf("courses/teacher_avatar/%s", courseID))
	if err != nil {
		return dto.CreateCourseResponse{}, err
	}

	_, err = s.fileUtil.ValidateAndUploadFile(ctx, req.Thumbnail, fileutil.ImageContentTypes,
		fmt.Sprintf("courses/thumbnail/%s", courseID))
	if err != nil {
		return dto.CreateCourseResponse{}, err
	}

	previewVideoPath := fmt.Sprintf("courses/preview_video/%s", courseID)
	previewVideoUploadURL, err := s.fileUtil.GetUploadSignedURL(previewVideoPath, "video/mp4")
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"path":  previewVideoPath,
		}, "Failed to get preview video signed upload URL")
		return dto.CreateCourseResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	course := &entity.Course{
		ID:          courseID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		TeacherName: req.TeacherName,
	}

	// Create course
	err = s.repo.CreateCourse(ctx, course)
	if err != nil {
		if strings.HasPrefix(err.Error(), "category not found") {
			return dto.CreateCourseResponse{}, errorpkg.ErrValidation().WithDetail("Category not found")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to create course")
		return dto.CreateCourseResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"course": course,
	}, "Course created")

	return dto.CreateCourseResponse{
		Course:                &dto.CourseResponse{ID: courseID},
		PreviewVideoUploadURL: previewVideoUploadURL,
	}, nil
}

func (s *courseService) GetCourseByID(ctx context.Context, id uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.repo.GetCourseByID(ctx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "course not found") {
			return nil, errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to get course by ID")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp := &dto.CourseResponse{}
	err = resp.PopulateFromEntity(course, s.fileUtil.GetSignedURL)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":  err,
			"course": course,
		}, "Failed to populate course response from entity")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return resp, nil
}

func (s *courseService) GetCourses(ctx context.Context, query dto.GetCoursesQuery,
	pageReq dto.PaginationRequest) ([]*dto.CourseResponse, dto.PaginationResponse, error) {
	courses, pageResp, err := s.repo.GetCourses(ctx, query, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"query": query,
			"page":  pageReq,
		}, "Failed to get courses")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp := make([]*dto.CourseResponse, len(courses))
	for i, course := range courses {
		resp[i] = &dto.CourseResponse{}
		err = resp[i].PopulateFromEntity(course, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":  err,
				"course": course,
			}, "Failed to populate course response from entity")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
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

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	err = s.repo.UpdateCourse(ctx, tx, updates)
	if err != nil {
		if strings.HasPrefix(err.Error(), "course not found") {
			return errorpkg.ErrNotFound()
		} else if strings.HasPrefix(err.Error(), "category not found") {
			return errorpkg.ErrValidation().WithDetail("Category not found")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to update course")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if req.TeacherAvatar != nil {
		_, err := s.fileUtil.ValidateAndUploadFile(ctx, req.TeacherAvatar, fileutil.ImageContentTypes,
			fmt.Sprintf("courses/teacher_avatar/%s", id))
		if err != nil {
			return err
		}
	}

	if req.Thumbnail != nil {
		_, err := s.fileUtil.ValidateAndUploadFile(ctx, req.Thumbnail, fileutil.ImageContentTypes,
			fmt.Sprintf("courses/thumbnail/%s", id))
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to commit transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"course.id": id,
		"request":   req,
	}, "Course updated")

	return nil
}

func (s *courseService) DeleteCourse(ctx context.Context, id uuid.UUID) error {
	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	err = s.repo.DeleteCourse(ctx, tx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "course not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to delete course")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("courses/teacher_avatar/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to delete teacher avatar")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("courses/thumbnail/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to delete thumbnail")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	err = s.fileUtil.Delete(ctx, fmt.Sprintf("courses/preview_video/%s", id))
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to delete preview video")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	err = tx.Commit()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to commit transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"course.id": id,
	}, "Course deleted")

	return nil
}

func (s *courseService) GetPreviewVideoUploadURL(ctx context.Context, id uuid.UUID) (string, error) {
	url, err := s.fileUtil.GetUploadSignedURL(fmt.Sprintf("courses/preview_video/%s", id), "video/mp4")
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"id":    id,
		}, "Failed to get preview video signed upload URL")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return url, nil
}

func (s *courseService) CreateEnrollment(ctx context.Context, courseID, studentID uuid.UUID) error {
	isSubscribedBoost, ok := ctx.Value(ctxkey.IsSubscribedBoost).(bool)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx, nil, "Failed to get isSubscribedBoost from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	if !isSubscribedBoost {
		return errorpkg.ErrForbiddenUser().WithDetail("You are not subscribed to Skill Boost")
	}

	err := s.repo.CreateEnrollment(ctx, courseID, studentID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "course not found") {
			return errorpkg.ErrValidation().WithDetail("Course not found")
		}
		if strings.HasPrefix(err.Error(), "student already enrolled in course") {
			return errorpkg.ErrStudentAlreadyEnrolled()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"course.id": courseID,
		}, "Failed to create enrollment")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"course.id": courseID,
	}, "Successfully created enrollment")

	return nil
}

func (s *courseService) GetEnrolledCourses(ctx context.Context, studentID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*dto.CourseResponse, dto.PaginationResponse, error) {
	courses, pageResp, err := s.repo.GetEnrolledCourses(ctx, studentID, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"page":  pageReq,
		}, "Failed to get enrolled courses")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp := make([]*dto.CourseResponse, len(courses))
	for i, course := range courses {
		resp[i] = &dto.CourseResponse{}
		err = resp[i].PopulateFromEntity(course, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":  err,
				"course": course,
			}, "Failed to populate course response from entity")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
		resp[i].PopulateFromCourseEnrollment(course.Enrollment)
	}

	return resp, pageResp, nil
}
