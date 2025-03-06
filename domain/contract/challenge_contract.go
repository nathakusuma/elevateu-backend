package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type IChallengeRepository interface {
	CreateChallenge(ctx context.Context, challenge *entity.Challenge) error
	GetChallenges(ctx context.Context, groupID uuid.UUID, difficulty enum.ChallengeDifficulty,
		pageReq dto.PaginationRequest) ([]*entity.Challenge, dto.PaginationResponse, error)
	GetChallengeByID(ctx context.Context, id uuid.UUID) (*entity.Challenge, error)
	UpdateChallenge(ctx context.Context, id uuid.UUID, updates *dto.ChallengeUpdate) error
	DeleteChallenge(ctx context.Context, id uuid.UUID) error
}

type IChallengeService interface {
	CreateChallenge(ctx context.Context, req *dto.CreateChallengeRequest) (*dto.ChallengeResponse, error)
	GetChallenges(ctx context.Context, groupID uuid.UUID, difficulty enum.ChallengeDifficulty,
		paginationReq dto.PaginationRequest) ([]*dto.ChallengeResponse, dto.PaginationResponse, error)
	GetChallengeDetail(ctx context.Context, id uuid.UUID) (*dto.ChallengeResponse, error)
	UpdateChallenge(ctx context.Context, id uuid.UUID, req *dto.UpdateChallengeRequest) error
	DeleteChallenge(ctx context.Context, id uuid.UUID) error
}
