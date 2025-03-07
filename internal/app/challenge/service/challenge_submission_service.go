package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type challengeSubmissionService struct {
	repo          contract.IChallengeSubmissionRepository
	challengeRepo contract.IChallengeRepository
	userRepo      contract.IUserRepository
	txManager     database.ITransactionManager
	fileUtil      fileutil.IFileUtil
	uuid          uuidpkg.IUUID
}

func NewChallengeSubmissionService(
	repo contract.IChallengeSubmissionRepository,
	challengeRepo contract.IChallengeRepository,
	userRepo contract.IUserRepository,
	txManager database.ITransactionManager,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.IChallengeSubmissionService {
	return &challengeSubmissionService{
		repo:          repo,
		challengeRepo: challengeRepo,
		userRepo:      userRepo,
		txManager:     txManager,
		fileUtil:      fileUtil,
		uuid:          uuid,
	}
}

func (s *challengeSubmissionService) CreateSubmission(ctx context.Context,
	req dto.CreateChallengeSubmissionRequest) error {
	userID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	challenge, err := s.challengeRepo.GetChallengeByID(ctx, req.ChallengeID)
	if err != nil {
		if err.Error() == "challenge not found" {
			return errorpkg.ErrValidation.Build().WithDetail("Challenge not found")
		}
	}

	submissionID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
			"user.id": userID,
		}, "[ChallengeSubmissionService][CreateSubmission] Failed to generate submission ID")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	submission := &entity.ChallengeSubmission{
		ID:          submissionID,
		ChallengeID: req.ChallengeID,
		StudentID:   userID,
		URL:         req.URL,
	}

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"submission": submission,
			"user.id":    userID,
		}, "[ChallengeSubmissionService][CreateSubmission] Failed to begin transaction for adding points")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}
	defer tx.Rollback()

	err = s.repo.CreateSubmission(ctx, tx, submission)
	if err != nil {
		if strings.Contains(err.Error(), "student has already submitted") {
			return errorpkg.ErrStudentAlreadySubmittedChallenge
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"submission": submission,
		}, "[ChallengeSubmissionService][CreateSubmission] Failed to create submission")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	var points int
	switch challenge.Difficulty {
	case enum.ChallengeDifficultyBeginner:
		points = 20
	case enum.ChallengeDifficultyIntermediate:
		points = 40
	case enum.ChallengeDifficultyAdvanced:
		points = 80
	}

	err = s.userRepo.AddPoint(ctx, tx, submission.StudentID, points)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"submission": submission,
			"points":     points,
		}, "[ChallengeSubmissionService][CreateSubmission] Failed to add points to student")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"submission": submission,
		}, "[ChallengeSubmissionService][CreateSubmission] Failed to commit transaction for adding points")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return nil
}

func (s *challengeSubmissionService) GetSubmissionAsStudent(ctx context.Context,
	challengeID uuid.UUID) (*dto.ChallengeSubmissionResponse, error) {
	studentID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return nil, errorpkg.ErrInvalidBearerToken
	}

	submission, err := s.repo.GetSubmissionByStudent(ctx, challengeID, studentID)
	if err != nil {
		if err.Error() == "submission not found" {
			return nil, errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"challengeID": challengeID,
			"studentID":   studentID,
		}, "[ChallengeSubmissionService][GetSubmissionAsStudent] Failed to get submission")
		return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	resp := &dto.ChallengeSubmissionResponse{}
	if err = resp.PopulateFromEntity(submission, s.fileUtil.GetSignedURL); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"submission": submission,
		}, "[ChallengeSubmissionService][GetSubmissionAsStudent] Failed to populate response")
		return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return resp, nil
}

func (s *challengeSubmissionService) GetSubmissionsAsMentor(ctx context.Context, challengeID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*dto.ChallengeSubmissionResponse, dto.PaginationResponse, error) {
	submissions, pageResp, err := s.repo.GetSubmissionsByChallenge(ctx, challengeID, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":       err,
			"challengeID": challengeID,
			"pagination":  pageReq,
		}, "[ChallengeSubmissionService][GetSubmissionsAsMentor] Failed to get submissions")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	resp := make([]*dto.ChallengeSubmissionResponse, len(submissions))
	for i, submission := range submissions {
		resp[i] = &dto.ChallengeSubmissionResponse{}
		if err = resp[i].PopulateFromEntity(submission, s.fileUtil.GetSignedURL); err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error":      err,
				"submission": submission,
			}, "[ChallengeSubmissionService][GetSubmissionsAsMentor] Failed to populate response")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
		}
	}

	return resp, pageResp, nil
}

func (s *challengeSubmissionService) CreateFeedback(ctx context.Context,
	submissionID uuid.UUID, req dto.CreateChallengeSubmissionFeedbackRequest) error {
	mentorID, ok := ctx.Value(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	submission, err := s.repo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		if err.Error() == "submission not found" {
			return errorpkg.ErrValidation.WithDetail("Submission not found")
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":         err,
			"submission.id": submissionID,
			"mentor.id":     mentorID,
		}, "[ChallengeSubmissionService][CreateFeedback] Failed to get submission")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	feedback := &entity.ChallengeSubmissionFeedback{
		SubmissionID: submissionID,
		MentorID:     mentorID,
		Score:        *req.Score,
		Feedback:     req.Feedback,
	}

	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"feedback":  feedback,
			"mentor.id": mentorID,
		}, "[ChallengeSubmissionService][CreateFeedback] Failed to begin transaction for adding points")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}
	defer tx.Rollback()

	err = s.repo.CreateFeedback(ctx, tx, feedback)
	if err != nil {
		if strings.Contains(err.Error(), "feedback already exists") {
			return errorpkg.ErrMentorAlreadySubmittedFeedback
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":     err,
			"feedback":  feedback,
			"mentor.id": mentorID,
		}, "[ChallengeSubmissionService][CreateFeedback] Failed to create feedback")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	var points int
	score := *req.Score
	switch {
	case score >= 85:
		points = 30
	case score >= 70:
		points = 20
	case score >= 50:
		points = 10
	}

	err = s.userRepo.AddPoint(ctx, tx, submission.StudentID, points)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"feedback": feedback,
			"points":   points,
		}, "[ChallengeSubmissionService][CreateFeedback] Failed to add points to student")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if err = tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":    err,
			"feedback": feedback,
		}, "[ChallengeSubmissionService][CreateFeedback] Failed to commit transaction for adding points")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return nil
}
