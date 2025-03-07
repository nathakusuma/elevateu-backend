package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type IChallengeSubmissionRepository interface {
	CreateSubmission(ctx context.Context, txWrapper database.ITransaction, submission *entity.ChallengeSubmission) error
	GetSubmissionByID(ctx context.Context, id uuid.UUID) (*entity.ChallengeSubmission, error)
	GetSubmissionByStudent(ctx context.Context, challengeID, studentID uuid.UUID) (*entity.ChallengeSubmission, error)
	GetSubmissionsByChallenge(ctx context.Context, challengeID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*entity.ChallengeSubmission, dto.PaginationResponse, error)

	CreateFeedback(ctx context.Context, txWrapper database.ITransaction,
		feedback *entity.ChallengeSubmissionFeedback) error
}

type IChallengeSubmissionService interface {
	CreateSubmission(ctx context.Context, req dto.CreateChallengeSubmissionRequest) error
	GetSubmissionAsStudent(ctx context.Context, challengeID uuid.UUID) (*dto.ChallengeSubmissionResponse, error)
	GetSubmissionsAsMentor(ctx context.Context, challengeID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*dto.ChallengeSubmissionResponse, dto.PaginationResponse, error)

	CreateFeedback(ctx context.Context, submissionID uuid.UUID, req dto.CreateChallengeSubmissionFeedbackRequest) error
}
