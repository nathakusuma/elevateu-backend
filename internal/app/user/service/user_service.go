package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/bcrypt"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type userService struct {
	repo     contract.IUserRepository
	bcrypt   bcrypt.IBcrypt
	fileUtil fileutil.IFileUtil
	uuid     uuidpkg.IUUID
}

func NewUserService(
	userRepo contract.IUserRepository,
	bcrypt bcrypt.IBcrypt,
	fileUtil fileutil.IFileUtil,
	uuid uuidpkg.IUUID,
) contract.IUserService {
	return &userService{
		repo:     userRepo,
		bcrypt:   bcrypt,
		fileUtil: fileUtil,
		uuid:     uuid,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (uuid.UUID, error) {
	// generate user ID
	userID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to generate user ID")

		return uuid.Nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	hash, err2 := s.bcrypt.Hash(req.Password)
	if err2 != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err2,
			"request": req,
		}, "Failed to hash password")

		return uuid.Nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// create user data
	user := &entity.User{
		ID:           userID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
	}

	// Add role-specific data
	if req.Role == enum.UserRoleStudent {
		if req.Student == nil {
			return uuid.Nil, errorpkg.ErrValidation().WithDetail("Student data is required")
		}
		user.Student = &entity.Student{
			Instance: req.Student.Instance,
			Major:    req.Student.Major,
		}
	} else if req.Role == enum.UserRoleMentor {
		if req.Mentor == nil {
			return uuid.Nil, errorpkg.ErrValidation().WithDetail("Mentor data is required")
		}
		user.Mentor = &entity.Mentor{
			Address:        req.Mentor.Address,
			Specialization: req.Mentor.Specialization,
			CurrentJob:     req.Mentor.CurrentJob,
			Company:        req.Mentor.Company,
			Gender:         req.Mentor.Gender,
		}
	}

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		// if email already exists
		if strings.HasPrefix(err.Error(), "conflict email") {
			return uuid.Nil, errorpkg.ErrEmailAlreadyRegistered()
		}

		// other error
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to create user")
		return uuid.Nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"user": user,
	}, "User created")

	return userID, nil
}

func (s *userService) getUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error) {
	user, err := s.repo.GetUserByField(ctx, field, value)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return nil, errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"field": field,
			"value": value,
		}, "Failed to get user by field")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.getUserByField(ctx, "email", email)
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID, isMinimal bool) (*dto.UserResponse, error) {
	user, err := s.getUserByField(ctx, "id", id)
	if err != nil {
		return nil, err
	}

	resp := &dto.UserResponse{}
	if isMinimal {
		err = resp.PopulateMinimalFromEntity(user, s.fileUtil.GetSignedURL)
	} else {
		err = resp.PopulateFromEntity(user, s.fileUtil.GetSignedURL)
	}

	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"user":  user,
		}, "Failed to populate user response")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return resp, nil
}

func (s *userService) UpdatePassword(ctx context.Context, email, newPassword string) error {
	// get user by email
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// hash new password
	newPasswordHash, err := s.bcrypt.Hash(newPassword)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "Failed to hash password")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	userUpdates := &dto.UserUpdate{
		ID:           user.ID,
		PasswordHash: &newPasswordHash,
	}

	if err = s.repo.UpdateUser(ctx, userUpdates); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "Failed to update user password")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"user.email": email,
	}, "Password updated")

	return nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) error {
	userUpdate := &dto.UserUpdate{
		ID:   id,
		Name: req.Name,
	}

	if req.Student != nil {
		userUpdate.Student = &dto.StudentUpdate{
			Instance: req.Student.Instance,
			Major:    req.Student.Major,
		}
	}

	if req.Mentor != nil {
		userUpdate.Mentor = &dto.MentorUpdate{
			Address:        req.Mentor.Address,
			Specialization: req.Mentor.Specialization,
			CurrentJob:     req.Mentor.CurrentJob,
			Company:        req.Mentor.Company,
			Bio:            req.Mentor.Bio,
			Gender:         req.Mentor.Gender,
		}
	}

	// update user
	if err := s.repo.UpdateUser(ctx, userUpdate); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "Failed to update user")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"updates": userUpdate,
	}, "User updated")

	return nil
}

func (s *userService) UpdateUserAvatar(ctx context.Context, id uuid.UUID, avatar *multipart.FileHeader) error {
	// get user by ID
	_, err := s.repo.GetUserByField(ctx, "id", id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"user.id": id,
		}, "Failed to get user")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// handle avatar upload
	_, err = s.fileUtil.ValidateAndUploadFile(ctx, avatar, fileutil.ImageContentTypes,
		fmt.Sprintf("users/avatar/%s", id.String()))
	if err != nil {
		return err
	}

	// update avatar URL
	hasAvatar := true
	userUpdate := &dto.UserUpdate{
		ID:        id,
		HasAvatar: &hasAvatar,
	}

	if err = s.repo.UpdateUser(ctx, userUpdate); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "Failed to update user avatar")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"user.id": id,
	}, "Avatar updated")

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// delete user
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to delete user")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// delete avatar
	err = s.fileUtil.Delete(ctx, fmt.Sprintf("users/avatar/%s", id.String()))
	if err != nil {
		if strings.Contains(err.Error(), "object doesn't exist") {
			goto pass
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to delete avatar")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

pass:
	log.Info(ctx, map[string]interface{}{
		"user.id": id,
	}, "User deleted")

	return nil
}

func (s *userService) DeleteUserAvatar(ctx context.Context, id uuid.UUID) error {
	// delete avatar
	if err := s.fileUtil.Delete(ctx, fmt.Sprintf("users/avatar/%s", id.String())); err != nil {
		if strings.Contains(err.Error(), "object doesn't exist") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to delete avatar")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// update avatar URL
	hasAvatar := false
	userUpdate := &dto.UserUpdate{
		ID:        id,
		HasAvatar: &hasAvatar,
	}

	if err := s.repo.UpdateUser(ctx, userUpdate); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "Failed to update user avatar")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, nil, "Avatar deleted")

	return nil
}

func (s *userService) GetLeaderboard(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := s.repo.GetTopPoints(ctx, 10)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to get top points")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	leaderboard := make([]*dto.UserResponse, len(users))
	for i, user := range users {
		leaderboard[i] = &dto.UserResponse{}
		if err = leaderboard[i].PopulateMinimalFromEntity(user, s.fileUtil.GetSignedURL); err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error": err,
				"user":  user,
			}, "Failed to populate user response")
			return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
	}

	return leaderboard, nil
}

func (s *userService) GetMentors(ctx context.Context, pageReq dto.PaginationRequest) ([]*dto.UserResponse, dto.PaginationResponse, error) {
	mentors, pageResp, err := s.repo.GetMentors(ctx, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"pagination": pageReq,
		}, "Failed to get mentors")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	responses := make([]*dto.UserResponse, len(mentors))
	for i, mentor := range mentors {
		responses[i] = &dto.UserResponse{}
		if err = responses[i].PopulateMinimalFromEntity(mentor, s.fileUtil.GetSignedURL); err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error": err,
				"user":  mentor,
			}, "Failed to populate mentor response")
			return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
	}

	return responses, pageResp, nil
}
