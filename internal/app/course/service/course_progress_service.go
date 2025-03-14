package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

type courseProgressService struct {
	repo      contract.ICourseProgressRepository
	txManager database.ITransactionManager
	userRepo  contract.IUserRepository
}

func NewCourseProgressService(
	progressRepo contract.ICourseProgressRepository,
	txManager database.ITransactionManager,
	userRepo contract.IUserRepository,
) contract.ICourseProgressService {
	return &courseProgressService{
		repo:      progressRepo,
		txManager: txManager,
		userRepo:  userRepo,
	}
}

func (s *courseProgressService) UpdateVideoProgress(ctx context.Context, studentID, videoID uuid.UUID,
	req dto.UpdateCourseVideoProgressRequest) error {
	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	courseID, err := s.repo.GetContentCourseID(ctx, videoID, "video")
	if err != nil {
		if strings.HasPrefix(err.Error(), "course content not found") {
			return errorpkg.ErrValidation().WithDetail("Video not found")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":    err,
			"video.id": videoID,
		}, "Failed to get course ID for video")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	progress := entity.CourseVideoProgress{
		StudentID:    studentID,
		VideoID:      videoID,
		LastPosition: req.LastPosition,
		IsCompleted:  req.IsCompleted,
	}

	newlyCompleted, err := s.repo.UpdateVideoProgress(ctx, tx, progress)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":    err,
			"video.id": videoID,
		}, "Failed to update video progress")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// If the video was newly completed, update the course progress
	if newlyCompleted {
		courseCompleted, err := s.repo.IncrementCourseProgress(ctx, tx, courseID, studentID)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":     err,
				"course.id": courseID,
			}, "Failed to update course progress")
			return errorpkg.ErrInternalServer().WithTraceID(traceID)
		}

		if courseCompleted {
			if err = s.userRepo.AddPoint(ctx, tx, studentID, 50); err != nil {
				traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
					"error":      err,
					"student.id": studentID,
				}, "Failed to add points to student")
				return errorpkg.ErrInternalServer().WithTraceID(traceID)
			}
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
		"video.id": videoID,
	}, "Video progress updated")

	return nil
}

func (s *courseProgressService) UpdateMaterialProgress(ctx context.Context, studentID uuid.UUID,
	materialID uuid.UUID) error {
	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	courseID, err := s.repo.GetContentCourseID(ctx, materialID, "material")
	if err != nil {
		if strings.HasPrefix(err.Error(), "course content not found") {
			return errorpkg.ErrValidation().WithDetail("Material not found")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":       err,
			"material.id": materialID,
		}, "Failed to get course ID for material")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	progress := entity.CourseMaterialProgress{
		StudentID:  studentID,
		MaterialID: materialID,
	}

	newlyCompleted, err := s.repo.UpdateMaterialProgress(ctx, tx, progress)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":       err,
			"material.id": materialID,
		}, "Failed to update material progress")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// If the material was newly completed, update the course progress
	if newlyCompleted {
		courseCompleted, err := s.repo.IncrementCourseProgress(ctx, tx, courseID, studentID)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":     err,
				"course.id": courseID,
			}, "Failed to update course progress")
			return errorpkg.ErrInternalServer().WithTraceID(traceID)
		}

		// If the course was just completed, add points to the student
		if courseCompleted {
			if err = s.userRepo.AddPoint(ctx, tx, studentID, 50); err != nil {
				traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
					"error":      err,
					"student.id": studentID,
				}, "Failed to add points to student")
				return errorpkg.ErrInternalServer().WithTraceID(traceID)
			}
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
		"material.id": materialID,
	}, "Material progress updated")

	return nil
}
