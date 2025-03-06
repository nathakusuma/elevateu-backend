package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type IChallengeGroupRepository interface {
	CreateGroup(ctx context.Context, group *entity.ChallengeGroup) error
	GetGroups(ctx context.Context, query dto.GetChallengeGroupQuery,
		pageReq dto.PaginationRequest) ([]*entity.ChallengeGroup, dto.PaginationResponse, error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, updates dto.ChallengeGroupUpdate) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
}

type IChallengeGroupService interface {
	CreateGroup(ctx context.Context, req dto.CreateChallengeGroupRequest) (*dto.ChallengeGroupResponse, error)
	GetGroups(ctx context.Context, query dto.GetChallengeGroupQuery,
		pageReq dto.PaginationRequest) ([]*dto.ChallengeGroupResponse, dto.PaginationResponse, error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, req dto.UpdateChallengeGroupRequest) error
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error
}
