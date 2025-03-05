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
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type courseFeedbackService struct {
	repo       contract.ICourseFeedbackRepository
	courseRepo contract.ICourseRepository
	fileUtil   fileutil.IFileUtil
	uuid       uuidpkg.IUUID
}

func NewCourseFeedbackService(
	courseFeedbackRepo contract.ICourseFeedbackRepository,
	courseRepo contract.ICourseRepository,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.ICourseFeedbackService {
	return &courseFeedbackService{
		repo:       courseFeedbackRepo,
		courseRepo: courseRepo,
		fileUtil:   fileUtil,
		uuid:       uuid,
	}
}

func (s *courseFeedbackService) CreateFeedback(ctx context.Context, courseID uuid.UUID,
	req dto.CreateCourseFeedbackRequest) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][CreateFeedback] Failed to begin transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}
	defer tx.Rollback()

	// Check if student has completed the course
	enrollment, err := s.courseRepo.GetEnrollment(ctx, courseID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "enrollment not found") {
			return errorpkg.ErrCannotFeedbackUnenrolledCourse
		}
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"course.id":  courseID,
			"student.id": userID,
		}, "[CourseFeedbackService][CreateFeedback] Failed to check if student completed course")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if !enrollment.IsCompleted {
		return errorpkg.ErrCannotFeedbackUncompletedCourse
	}

	feedbackID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][CreateFeedback] Failed to generate UUID")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
		if strings.Contains(err.Error(), "student has already submitted feedback") {
			return errorpkg.ErrStudentAlreadySubmittedFeedback
		}

		if strings.Contains(err.Error(), "course not found") {
			return errorpkg.ErrValidation.Build().WithDetail("Course not found")
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"feedback":  feedback,
			"course.id": courseID,
		}, "[CourseFeedbackService][CreateFeedback] Failed to create feedback")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	course, err := s.courseRepo.GetCourseByID(ctx, courseID)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"course.id": courseID,
		}, "[CourseFeedbackService][CreateFeedback] Failed to get course")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}
	ratingCount := course.RatingCount
	totalRating := course.TotalRating

	newRatingCount := ratingCount + 1
	newTotalRating := totalRating + float64(req.Rating)
	newRating := newTotalRating / float64(newRatingCount)

	// Update course with new rating
	err = s.repo.UpdateCourseRating(ctx, tx, courseID, newRatingCount, newRating, newTotalRating)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"course.id": courseID,
		}, "[CourseFeedbackService][CreateFeedback] Failed to update course rating")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][CreateFeedback] Failed to commit transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"feedback": feedback,
	}, "[CourseFeedbackService][CreateFeedback] Feedback created successfully")

	return nil
}

func (s *courseFeedbackService) GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*dto.CourseFeedbackResponse, dto.PaginationResponse, error) {

	feedbacks, pageResp, err := s.repo.GetFeedbacksByCourseID(ctx, courseID, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"course.id":  courseID,
			"pagination": pageReq,
		}, "[CourseFeedbackService][GetFeedbacksByCourseID] Failed to get feedbacks")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	responses := make([]*dto.CourseFeedbackResponse, len(feedbacks))
	for i, feedback := range feedbacks {
		responses[i] = &dto.CourseFeedbackResponse{}
		err = responses[i].PopulateFromEntity(feedback, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":    err,
				"feedback": feedback,
			}, "[CourseFeedbackService][GetFeedbacksByCourseID] Failed to populate response from entity")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
		}
	}

	return responses, pageResp, nil
}

func (s *courseFeedbackService) UpdateFeedback(ctx context.Context, feedbackID uuid.UUID,
	req dto.UpdateCourseFeedbackRequest) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	feedback, err := s.repo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		if strings.Contains(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
		}, "[CourseFeedbackService][UpdateFeedback] Failed to get feedback")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if feedback.StudentID != userID {
		return errorpkg.ErrForbiddenUser.WithDetail("You can only update your own feedback")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][UpdateFeedback] Failed to begin transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
			"updates":     updates,
		}, "[CourseFeedbackService][UpdateFeedback] Failed to update feedback")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if ratingDiff != 0 {
		course, err := s.courseRepo.GetCourseByID(ctx, feedback.CourseID)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":     err,
				"course.id": feedback.CourseID,
			}, "[CourseFeedbackService][UpdateFeedback] Failed to get course")
			return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
		}
		ratingCount := course.RatingCount
		totalRating := course.TotalRating

		newTotalRating := totalRating + ratingDiff
		newRating := newTotalRating / float64(ratingCount)

		err = s.repo.UpdateCourseRating(ctx, tx, feedback.CourseID, ratingCount, newRating, newTotalRating)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":     err,
				"course.id": feedback.CourseID,
			}, "[CourseFeedbackService][UpdateFeedback] Failed to update course rating")
			return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
		}
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][UpdateFeedback] Failed to commit transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"feedback.id": feedbackID,
		"student.id":  userID,
		"course.id":   feedback.CourseID,
		"updates":     updates,
	}, "[CourseFeedbackService][UpdateFeedback] Feedback updated successfully")

	return nil
}

func (s *courseFeedbackService) DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	feedback, err := s.repo.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		if strings.Contains(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
		}, "[CourseFeedbackService][DeleteFeedback] Failed to get feedback")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if feedback.StudentID != userID {
		return errorpkg.ErrForbiddenUser.WithDetail("You can only delete your own feedback")
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][DeleteFeedback] Failed to begin transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	defer tx.Rollback()

	err = s.repo.DeleteFeedback(ctx, tx, feedbackID)
	if err != nil {
		if strings.Contains(err.Error(), "feedback not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"feedback.id": feedbackID,
		}, "[CourseFeedbackService][DeleteFeedback] Failed to delete feedback")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	course, err := s.courseRepo.GetCourseByID(ctx, feedback.CourseID)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"course.id": feedback.CourseID,
		}, "[CourseFeedbackService][DeleteFeedback] Failed to get course")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"course.id": feedback.CourseID,
		}, "[CourseFeedbackService][DeleteFeedback] Failed to update course rating")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
		}, "[CourseFeedbackService][DeleteFeedback] Failed to commit transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"feedback.id": feedbackID,
		"course.id":   feedback.CourseID,
		"student.id":  userID,
	}, "[CourseFeedbackService][DeleteFeedback] Feedback deleted successfully")

	return nil
}
