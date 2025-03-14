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

type challengeGroupService struct {
	repo     contract.IChallengeGroupRepository
	fileUtil fileutil.IFileUtil
	uuid     uuidpkg.IUUID
}

func NewChallengeGroupService(
	repo contract.IChallengeGroupRepository,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.IChallengeGroupService {
	return &challengeGroupService{
		repo:     repo,
		fileUtil: fileUtil,
		uuid:     uuid,
	}
}

func (s *challengeGroupService) CreateGroup(ctx context.Context,
	req dto.CreateChallengeGroupRequest) (*dto.ChallengeGroupResponse, error) {
	groupID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[ChallengeGroupService][CreateGroup] Failed to generate group ID")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if req.Thumbnail != nil {
		_, err = s.fileUtil.ValidateAndUploadFile(ctx, req.Thumbnail, fileutil.ImageContentTypes,
			fmt.Sprintf("challenge_groups/thumbnail/%s", groupID))
		if err != nil {
			return nil, err
		}
	}

	group := &entity.ChallengeGroup{
		ID:          groupID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
	}

	err = s.repo.CreateGroup(ctx, group)
	if err != nil {
		if strings.HasPrefix(err.Error(), "category not found") {
			return nil, errorpkg.ErrValidation().WithDetail("Category not found")
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[ChallengeGroupService][CreateGroup] Failed to create challenge group")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return &dto.ChallengeGroupResponse{
		ID: groupID,
	}, nil
}

func (s *challengeGroupService) GetGroups(ctx context.Context, query dto.GetChallengeGroupQuery,
	pageReq dto.PaginationRequest) ([]*dto.ChallengeGroupResponse, dto.PaginationResponse, error) {
	groups, pageResp, err := s.repo.GetGroups(ctx, query, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"query": query,
			"page":  pageReq,
		}, "[ChallengeGroupService][GetGroups] Failed to get challenge groups")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp := make([]*dto.ChallengeGroupResponse, len(groups))
	for i, group := range groups {
		resp[i] = &dto.ChallengeGroupResponse{}
		err = resp[i].PopulateFromEntity(group, s.fileUtil.GetSignedURL)
		if err != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error": err,
				"group": group,
			}, "[ChallengeGroupService][GetGroups] Failed to populate response")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
	}

	return resp, pageResp, nil
}

func (s *challengeGroupService) UpdateGroup(ctx context.Context, groupID uuid.UUID,
	req dto.UpdateChallengeGroupRequest) error {
	updates := dto.ChallengeGroupUpdate{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
	}

	err := s.repo.UpdateGroup(ctx, groupID, updates)
	if err != nil {
		if err.Error() == "challenge group not found" {
			return errorpkg.ErrNotFound()
		} else if strings.HasPrefix(err.Error(), "category not found") {
			return errorpkg.ErrValidation().WithDetail("Category not found")
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"groupID": groupID,
			"request": req,
		}, "[ChallengeGroupService][UpdateGroup] Failed to update challenge group")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if req.Thumbnail != nil {
		_, err := s.fileUtil.ValidateAndUploadFile(ctx, req.Thumbnail, fileutil.ImageContentTypes,
			fmt.Sprintf("challenge_groups/thumbnail/%s", groupID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *challengeGroupService) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	err := s.repo.DeleteGroup(ctx, groupID)
	if err != nil {
		if err.Error() == "challenge group not found" {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"groupID": groupID,
		}, "[ChallengeGroupService][DeleteGroup] Failed to delete challenge group")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// Delete thumbnail file
	err = s.fileUtil.Delete(ctx, fmt.Sprintf("challenge_groups/thumbnail/%s", groupID))
	if err != nil {
		log.Error(map[string]interface{}{
			"error":   err,
			"groupID": groupID,
		}, "[ChallengeGroupService][DeleteGroup] Failed to delete thumbnail")
	}

	return nil
}
