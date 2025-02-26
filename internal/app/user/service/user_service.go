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
	userRepo contract.IUserRepository
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
		userRepo: userRepo,
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

		return uuid.Nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	hash, err2 := s.bcrypt.Hash(req.Password)
	if err2 != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":        err2,
			"request":      req,
			"requester.id": creatorID,
		}, "[UserService][CreateUser] Failed to hash password")

		return uuid.Nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
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
			return uuid.Nil, errorpkg.ErrValidation.WithDetail("Student data is required")
		}
		user.Student = &entity.Student{
			Instance: req.Student.Instance,
			Major:    req.Student.Major,
		}
	} else if req.Role == enum.UserRoleMentor {
		if req.Mentor == nil {
			return uuid.Nil, errorpkg.ErrValidation.WithDetail("Mentor data is required")
		}
		user.Mentor = &entity.Mentor{
			Specialization: req.Mentor.Specialization,
			Experience:     req.Mentor.Experience,
			Price:          req.Mentor.Price,
		}
	}

	err = s.userRepo.CreateUser(ctx, user)
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
		return uuid.Nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"user":         user,
		"requester.id": creatorID,
	}, "[UserService][CreateUser] User created")

	return userID, nil
}

func (s *userService) getUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error) {
	// get from repository
	user, err := s.userRepo.GetUserByField(ctx, field, value)
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
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	if user.AvatarURL != nil {
		avatarURL, err2 := s.fileUtil.GetSignedURL(*user.AvatarURL)
		if err2 != nil {
			traceID := log.ErrorWithTraceID(map[string]interface{}{
				"error": err,
				"user":  user,
			}, "[UserService][getUserByField] Failed to get avatar URL")
			return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
		}

		user.AvatarURL = &avatarURL
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.getUserByField(ctx, "email", email)
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.getUserByField(ctx, "id", id)
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
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	userUpdates := &dto.UserUpdate{
		ID:           user.ID,
		PasswordHash: &newPasswordHash,
	}

	if err = s.userRepo.UpdateUser(ctx, userUpdates); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "[UserService][UpdatePassword] Failed to update user password")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
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
			Specialization: req.Mentor.Specialization,
			Experience:     req.Mentor.Experience,
			Price:          req.Mentor.Price,
		}
	}

	if req.Avatar != nil {
		avatarURL, err2 := s.handleAvatarUpload(ctx, req.Avatar, id)
		if err2 != nil {
			return err2
		}
		userUpdate.AvatarURL = avatarURL
	}

	// update user
	if err := s.userRepo.UpdateUser(ctx, userUpdate); err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"updates": userUpdate,
		}, "[UserService][UpdateUser] Failed to update user")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"updates": userUpdate,
	}, "[UserService][UpdateUser] User updated")

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	requesterID := ctx.Value(ctxkey.UserID)
	if requesterID == nil {
		requesterID = "system"
	}

	// delete user
	err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":        err,
			"user.id":      id,
			"requester.id": requesterID,
		}, "[UserService][DeleteUser] Failed to delete user")
		return errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	log.Info(map[string]interface{}{
		"user.id":      id,
		"requester.id": requesterID,
	}, "[UserService][DeleteUser] User deleted")

	return nil
}

func (s *userService) handleAvatarUpload(ctx context.Context, avatar *multipart.FileHeader,
	userID uuid.UUID) (*string, error) {
	file, err := avatar.Open()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": userID,
		}, "[UserService][handleAvatarUpload] Failed to open avatar file")
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	defer file.Close()

	if avatar.Size > 2*fileutil.MegaByte {
		return nil, errorpkg.ErrFileTooLarge.WithDetail(
			fmt.Sprintf("File size is too large (%s). Please upload a file less than 2MB",
				fileutil.ByteToAppropriateUnit(avatar.Size)))
	}

	ok, fileType, err := s.fileUtil.CheckMIMEFileType(file, fileutil.ImageContentTypes)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": userID,
		}, "[UserService][handleAvatarUpload] Failed to check MIME file type")
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}
	if !ok {
		return nil, errorpkg.ErrInvalidFileFormat.WithDetail(
			fmt.Sprintf("File type %s is not allowed. Please upload a valid image file", fileType))
	}

	avatarURL, err := s.fileUtil.Upload(ctx, file, fmt.Sprintf("users/avatar/%s", userID.String()))
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"user.id": userID,
		}, "[UserService][handleAvatarUpload] Failed to upload avatar")
		return nil, errorpkg.ErrInternalServer.WithTraceID(traceID)
	}

	return &avatarURL, nil
}
