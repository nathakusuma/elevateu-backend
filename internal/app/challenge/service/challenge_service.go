package service

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type challengeService struct {
	repo     contract.IChallengeRepository
	fileUtil fileutil.IFileUtil
	uuid     uuidpkg.IUUID
}

func NewChallengeService(
	repo contract.IChallengeRepository,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.IChallengeService {
	return &challengeService{
		repo:     repo,
		fileUtil: fileUtil,
		uuid:     uuid,
	}
}

func (s *challengeService) CreateChallenge(ctx context.Context,
	req *dto.CreateChallengeRequest) (*dto.ChallengeResponse, error) {
	challengeID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[ChallengeService][CreateChallenge] Failed to generate challenge ID")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	challenge := &entity.Challenge{
		ID:          challengeID,
		GroupID:     req.GroupID,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		IsFree:      req.IsFree,
	}

	err = s.repo.CreateChallenge(ctx, challenge)
	if err != nil {
		if strings.HasPrefix(err.Error(), "challenge group not found") {
			return nil, errorpkg.ErrValidation().WithDetail("Challenge group not found")
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[ChallengeService][CreateChallenge] Failed to create challenge")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return &dto.ChallengeResponse{ID: challengeID}, nil
}

func (s *challengeService) GetChallenges(ctx context.Context, groupID uuid.UUID, difficulty enum.ChallengeDifficulty,
	paginationReq dto.PaginationRequest) ([]*dto.ChallengeResponse, dto.PaginationResponse, error) {
	challenges, pageResp, err := s.repo.GetChallenges(ctx, groupID, difficulty, paginationReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"groupID":    groupID,
			"difficulty": difficulty,
			"pagination": paginationReq,
		}, "[ChallengeService][GetChallenges] Failed to get challenges")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp := make([]*dto.ChallengeResponse, len(challenges))
	for i, challenge := range challenges {
		resp[i] = &dto.ChallengeResponse{}
		resp[i].PopulateFromEntity(challenge)
	}

	return resp, pageResp, nil
}

func (s *challengeService) GetChallengeDetail(ctx context.Context, id uuid.UUID) (*dto.ChallengeResponse, error) {
	challenge, err := s.repo.GetChallengeByID(ctx, id)
	if err != nil {
		if err.Error() == "challenge not found" {
			return nil, errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[ChallengeService][GetChallengeDetail] Failed to get challenge detail")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp := &dto.ChallengeResponse{}
	resp.PopulateDetailFromEntity(challenge)

	return resp, nil
}

func (s *challengeService) UpdateChallenge(ctx context.Context, id uuid.UUID, req *dto.UpdateChallengeRequest) error {
	updates := &dto.ChallengeUpdate{
		GroupID:     req.GroupID,
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		IsFree:      req.IsFree,
	}

	err := s.repo.UpdateChallenge(ctx, id, updates)
	if err != nil {
		if err.Error() == "challenge not found" {
			return errorpkg.ErrNotFound()
		} else if strings.HasPrefix(err.Error(), "challenge group not found") {
			return errorpkg.ErrValidation().WithDetail("Challenge group not found")
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"id":      id,
			"request": req,
		}, "[ChallengeService][UpdateChallenge] Failed to update challenge")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return nil
}

func (s *challengeService) DeleteChallenge(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteChallenge(ctx, id)
	if err != nil {
		if err.Error() == "challenge not found" {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"id":    id,
		}, "[ChallengeService][DeleteChallenge] Failed to delete challenge")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return nil
}
