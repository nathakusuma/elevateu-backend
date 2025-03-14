package service

import (
	"context"
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

type courseFeedbackService struct {
	repo       contract.ICourseFeedbackRepository
	courseRepo contract.ICourseRepository
	fileUtil   fileutil.IFileUtil
	txManager  database.ITransactionManager
	uuid       uuidpkg.IUUID
}

func NewCourseFeedbackService(
	courseFeedbackRepo contract.ICourseFeedbackRepository,
	courseRepo contract.ICourseRepository,
	fileUtil fileutil.IFileUtil,
	txManager database.ITransactionManager,
	uuid uuidpkg.IUUID,
) contract.ICourseFeedbackService {
	return &courseFeedbackService{
		repo:       courseFeedbackRepo,
		courseRepo: courseRepo,
		fileUtil:   fileUtil,
		txManager:  txManager,
		uuid:       uuid,
	}
}

func (s *courseFeedbackService) CreateFeedback(ctx context.Context, courseID uuid.UUID,
	req dto.CreateCourseFeedbackRequest) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx, nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][CreateFeedback] Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	// Check if student has completed the course
	enrollment, err := s.courseRepo.GetEnrollment(ctx, courseID, userID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "enrollment not found") {
			return errorpkg.ErrCannotFeedbackUnenrolledCourse()
		}
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"course.id": courseID,
		}, "Failed to check if student completed course")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if !enrollment.IsCompleted {
		return errorpkg.ErrCannotFeedbackUncompletedCourse()
	}

	feedbackID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to generate UUID")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	feedback := &entity.CourseFeedback{
		ID:        feedbackID,
		CourseID:  courseID,
		StudentID: userID,
		Rating:    float64(req.Rating),
		Comment:   req.Comment,
	}

	err = s.repo.CreateFeedback(ctx, tx, feedback)
	if err != nil {
		if strings.HasPrefix(err.Error(), "student has already submitted feedback for this course") {
			return errorpkg.ErrStudentAlreadySubmittedFeedback()
		}

		if strings.HasPrefix(err.Error(), "course not found") {
			return errorpkg.ErrValidation().WithDetail("Course not found")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"feedback":  feedback,
			"course.id": courseID,
		}, "Failed to create feedback")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	course, err := s.courseRepo.GetCourseByID(ctx, courseID)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"course.id": courseID,
		}, "Failed to get course")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	ratingCount := course.RatingCount
	totalRating := course.TotalRating

	newRatingCount := ratingCount + 1
	newTotalRating := totalRating + float64(req.Rating)
	newRating := newTotalRating / float64(newRatingCount)

	// Update course with new rating
	err = s.repo.UpdateCourseRating(ctx, tx, courseID, newRatingCount, newRating, newTotalRating)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"course.id": courseID,
		}, "Failed to update course rating")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to commit transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"feedback": feedback,
	}, "Feedback created successfully")

	return nil
}

func (s *courseFeedbackService) GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*dto.CourseFeedbackResponse, dto.PaginationResponse, error) {
	feedbacks, pageResp, err := s.repo.GetFeedbacksByCourseID(ctx, courseID, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"course.id":  courseID,
			"pagination": pageReq,
		}, "Failed to get feedbacks")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	responses := make([]*dto.CourseFeedbackResponse, len(feedbacks))
	for i, feedback := range feedbacks {
		responses[i] = &dto.CourseFeedbackResponse{}
		err = responses[i].PopulateFromEntity(feedback, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":    err,
				"feedback": feedback,
			}, "Failed to populate response from entity")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
	}

	return responses, pageResp, nil
}

func (s *courseFeedbackService) UpdateFeedback(ctx context.Context, feedbackID uuid.UUID,
	req dto.UpdateCourseFeedbackRequest) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx, nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	feedback, err := s.repo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		if strings.Contains(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
		}, "Failed to get feedback")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if feedback.StudentID != userID {
		return errorpkg.ErrForbiddenUser().WithDetail("You can only update your own feedback")
	}

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	defer tx.Rollback()

	updates := dto.CourseFeedbackUpdate{}

	var ratingDiff float64 = 0

	if req.Rating != 0 {
		newRating := float64(req.Rating)
		updates.Rating = &newRating
		ratingDiff = newRating - feedback.Rating
	}

	if req.Comment != "" {
		updates.Feedback = &req.Comment
	}

	err = s.repo.UpdateFeedback(ctx, tx, feedbackID, updates)
	if err != nil {
		if strings.Contains(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
			"updates":     updates,
		}, "Failed to update feedback")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if ratingDiff != 0 {
		course, err := s.courseRepo.GetCourseByID(ctx, feedback.CourseID)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":     err,
				"course.id": feedback.CourseID,
			}, "Failed to get course")
			return errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
		ratingCount := course.RatingCount
		totalRating := course.TotalRating

		newTotalRating := totalRating + ratingDiff
		newRating := newTotalRating / float64(ratingCount)

		err = s.repo.UpdateCourseRating(ctx, tx, feedback.CourseID, ratingCount, newRating, newTotalRating)
		if err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":     err,
				"course.id": feedback.CourseID,
			}, "Failed to update course rating")
			return errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to commit transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"feedback.id": feedbackID,
		"course.id":   feedback.CourseID,
		"updates":     updates,
	}, "Feedback updated successfully")

	return nil
}

func (s *courseFeedbackService) DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx, nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	feedback, err := s.repo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
		}, "Failed to get feedback")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if feedback.StudentID != userID {
		return errorpkg.ErrForbiddenUser().WithDetail("You can only delete your own feedback")
	}

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	defer tx.Rollback()

	err = s.repo.DeleteFeedback(ctx, tx, feedbackID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
		}, "Failed to delete feedback")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	course, err := s.courseRepo.GetCourseByID(ctx, feedback.CourseID)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"course.id": feedback.CourseID,
		}, "Failed to get course")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	ratingCount := course.RatingCount
	totalRating := course.TotalRating

	newRatingCount := ratingCount - 1
	newTotalRating := totalRating - feedback.Rating

	var newRating float64
	if newRatingCount > 0 {
		newRating = newTotalRating / float64(newRatingCount)
	}

	err = s.repo.UpdateCourseRating(ctx, tx, feedback.CourseID, newRatingCount, newRating, newTotalRating)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"course.id": feedback.CourseID,
		}, "Failed to update course rating")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to commit transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"feedback.id": feedbackID,
		"course.id":   feedback.CourseID,
	}, "Feedback deleted successfully")

	return nil
}
