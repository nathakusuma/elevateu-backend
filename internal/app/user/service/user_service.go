package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
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
	creatorID := ctx.Value(ctxkey.UserID)
	if creatorID == nil {
		creatorID = "system"
	}

	// generate user ID
	userID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":        err,
			"request":      req,
			"requester.id": creatorID,
		}, "[UserService][CreateUser] Failed to generate user ID")

		return uuid.Nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	hash, err2 := s.bcrypt.Hash(req.Password)
	if err2 != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":        err2,
			"request":      req,
			"requester.id": creatorID,
		}, "[UserService][CreateUser] Failed to hash password")

		return uuid.Nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
			return uuid.Nil, errorpkg.ErrValidation.Build().WithDetail("Student data is required")
		}
		user.Student = &entity.Student{
			Instance: req.Student.Instance,
			Major:    req.Student.Major,
		}
	} else if req.Role == enum.UserRoleMentor {
		if req.Mentor == nil {
			return uuid.Nil, errorpkg.ErrValidation.Build().WithDetail("Mentor data is required")
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
			return uuid.Nil, errorpkg.ErrEmailAlreadyRegistered
		}

		// other error
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":        err,
			"request":      req,
			"requester.id": creatorID,
		}, "[UserService][CreateUser] Failed to create user")
		return uuid.Nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"user":         user,
		"requester.id": creatorID,
	}, "[UserService][CreateUser] User created")

	return userID, nil
}

func (s *userService) getUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error) {
	// get from repository
	user, err := s.repo.GetUserByField(ctx, field, value)
	if err != nil {
		// if user not found
		if strings.HasPrefix(err.Error(), "user not found") {
			return nil, errorpkg.ErrNotFound
		}

		// other error
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"field": field,
			"value": value,
		}, "[UserService][getUserByField] Failed to get user by field")
		return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error": err,
			"user":  user,
		}, "[UserService][GetUserByID] Failed to populate user response")
		return nil, errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "[UserService][UpdatePassword] Failed to hash password")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	userUpdates := &dto.UserUpdate{
		ID:           user.ID,
		PasswordHash: &newPasswordHash,
	}

	if err = s.repo.UpdateUser(ctx, userUpdates); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "[UserService][UpdatePassword] Failed to update user password")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"user.email": email,
	}, "[UserService][UpdatePassword] Password updated")

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
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "[UserService][UpdateUser] Failed to update user")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"updates": userUpdate,
	}, "[UserService][UpdateUser] User updated")

	return nil
}

func (s *userService) UpdateUserAvatar(ctx context.Context, id uuid.UUID, avatar *multipart.FileHeader) error {
	// get user by ID
	_, err := s.repo.GetUserByField(ctx, "id", id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": id,
		}, "[UserService][UpdateUserAvatar] Failed to get user")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
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
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "[UserService][UpdateUserAvatar] Failed to update user avatar")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"user.id": id,
	}, "[UserService][UpdateUserAvatar] Avatar updated")

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// delete user
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": id,
		}, "[UserService][DeleteUser] Failed to delete user")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// delete avatar
	err = s.fileUtil.Delete(ctx, fmt.Sprintf("users/avatar/%s", id.String()))
	if err != nil {
		if strings.Contains(err.Error(), "object doesn't exist") {
			goto pass
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": id,
		}, "[UserService][DeleteUser] Failed to delete avatar")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

pass:
	log.Info(map[string]interface{}{
		"user.id": id,
	}, "[UserService][DeleteUser] User deleted")

	return nil
}

func (s *userService) DeleteUserAvatar(ctx context.Context, id uuid.UUID) error {
	// delete avatar
	if err := s.fileUtil.Delete(ctx, fmt.Sprintf("users/avatar/%s", id.String())); err != nil {
		if strings.Contains(err.Error(), "object doesn't exist") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": id,
		}, "[UserService][DeleteUserAvatar] Failed to delete avatar")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// update avatar URL
	hasAvatar := false
	userUpdate := &dto.UserUpdate{
		ID:        id,
		HasAvatar: &hasAvatar,
	}

	if err := s.repo.UpdateUser(ctx, userUpdate); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "[UserService][DeleteUserAvatar] Failed to update user avatar")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"user.id": id,
	}, "[UserService][DeleteUserAvatar] Avatar deleted")

	return nil
}
