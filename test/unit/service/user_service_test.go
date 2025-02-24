package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/app/user/service"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	appmocks "github.com/nathakusuma/elevateu-backend/test/unit/mocks/app"
	pkgmocks "github.com/nathakusuma/elevateu-backend/test/unit/mocks/pkg"
	_ "github.com/nathakusuma/elevateu-backend/test/unit/setup" // Initialize test environment
)

type userServiceMocks struct {
	userRepo    *appmocks.MockIUserRepository
	storageRepo *appmocks.MockIStorageRepository
	bcrypt      *pkgmocks.MockIBcrypt
	fileUtil    *pkgmocks.MockIFileUtil
	uuid        *pkgmocks.MockIUUID
}

func setupUserServiceTest(t *testing.T) (contract.IUserService, *userServiceMocks) {
	mocks := &userServiceMocks{
		userRepo:    appmocks.NewMockIUserRepository(t),
		storageRepo: appmocks.NewMockIStorageRepository(t),
		bcrypt:      pkgmocks.NewMockIBcrypt(t),
		fileUtil:    pkgmocks.NewMockIFileUtil(t),
		uuid:        pkgmocks.NewMockIUUID(t),
	}

	svc := service.NewUserService(mocks.userRepo, mocks.storageRepo, mocks.bcrypt, mocks.fileUtil, mocks.uuid)

	return svc, mocks
}

func Test_UserService_CreateUser(t *testing.T) {
	ctx := context.Background()
	hashedPassword := "hashed_password"
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := &dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: hashedPassword,
			Role:     enum.RoleStudent,
		}

		// Expect UUID generation
		mocks.uuid.EXPECT().
			NewV7().
			Return(userID, nil)

		// Expect password hashing
		mocks.bcrypt.EXPECT().
			Hash(req.Password).
			Return(hashedPassword, nil)

		// Expect user creation
		mocks.userRepo.EXPECT().
			CreateUser(ctx, &entity.User{
				ID:           userID,
				Name:         req.Name,
				Email:        req.Email,
				PasswordHash: req.Password,
				Role:         req.Role,
			}).
			Return(nil)

		resultID, err := svc.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, userID, resultID)
	})

	t.Run("error - password hashing fails", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := &dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		// Expect UUID generation
		mocks.uuid.EXPECT().
			NewV7().
			Return(userID, nil)

		// Expect password hashing to fail
		mocks.bcrypt.EXPECT().
			Hash(req.Password).
			Return("", errors.New("hashing error"))

		resultID, err := svc.CreateUser(ctx, req)
		assert.Equal(t, uuid.Nil, resultID)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - uuid generation fails", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := &dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		// Expect UUID generation to fail
		mocks.uuid.EXPECT().
			NewV7().
			Return(uuid.UUID{}, errors.New("uuid error"))

		resultID, err := svc.CreateUser(ctx, req)
		assert.Equal(t, uuid.Nil, resultID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - email already exists", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := &dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		// Expect UUID generation
		mocks.uuid.EXPECT().
			NewV7().
			Return(userID, nil)

		// Expect password hashing
		mocks.bcrypt.EXPECT().
			Hash(req.Password).
			Return(hashedPassword, nil)

		// Expect user creation to fail with conflict email error
		mocks.userRepo.EXPECT().
			CreateUser(ctx, &entity.User{
				ID:           userID,
				Name:         req.Name,
				Email:        req.Email,
				PasswordHash: req.Password,
				Role:         req.Role,
			}).
			Return(errors.New("conflict email: blablabla"))

		resultID, err := svc.CreateUser(ctx, req)
		assert.Equal(t, uuid.Nil, resultID)
		assert.ErrorIs(t, err, errorpkg.ErrEmailAlreadyRegistered)
	})

	t.Run("error - repository error", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := &dto.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: hashedPassword,
		}

		// Expect UUID generation
		mocks.uuid.EXPECT().
			NewV7().
			Return(userID, nil)

		hashedPassword2 := "hashed_password"
		mocks.bcrypt.EXPECT().
			Hash(req.Password).
			Return(hashedPassword2, nil)

		// Expect user creation to fail
		mocks.userRepo.EXPECT().
			CreateUser(ctx, &entity.User{
				ID:           userID,
				Name:         req.Name,
				Email:        req.Email,
				PasswordHash: hashedPassword,
				Role:         req.Role,
			}).
			Return(errors.New("db error"))

		resultID, err := svc.CreateUser(ctx, req)
		assert.Equal(t, uuid.Nil, resultID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_UserService_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"

	t.Run("success - without avatar", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		expectedUser := &entity.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     email,
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(expectedUser, nil)

		user, err := svc.GetUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("success - with avatar", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		avatarURL := "storage/avatar.jpg"
		signedURL := "https://storage.com/signed/avatar.jpg"

		expectedUser := &entity.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     email,
			Role:      enum.RoleStudent,
			AvatarURL: &avatarURL,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(expectedUser, nil)

		mocks.storageRepo.EXPECT().
			GetSignedURL(avatarURL).
			Return(signedURL, nil)

		user, err := svc.GetUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, signedURL, *user.AvatarURL)
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(nil, errors.New("user not found: blablabla"))

		user, err := svc.GetUserByEmail(ctx, email)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - repository error", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(nil, errors.New("db error"))

		user, err := svc.GetUserByEmail(ctx, email)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - storage signed URL error", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		avatarURL := "storage/avatar.jpg"
		expectedUser := &entity.User{
			ID:        uuid.New(),
			Name:      "Test User",
			Email:     email,
			Role:      enum.RoleStudent,
			AvatarURL: &avatarURL,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(expectedUser, nil)

		mocks.storageRepo.EXPECT().
			GetSignedURL(avatarURL).
			Return("", errors.New("storage error"))

		user, err := svc.GetUserByEmail(ctx, email)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_UserService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	t.Run("success - without avatar", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		expectedUser := &entity.User{
			ID:        id,
			Name:      "Test User",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", id.String()).
			Return(expectedUser, nil)

		user, err := svc.GetUserByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("success - with avatar", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		avatarURL := "storage/avatar.jpg"
		signedURL := "https://storage.com/signed/avatar.jpg"

		expectedUser := &entity.User{
			ID:        id,
			Name:      "Test User",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			AvatarURL: &avatarURL,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", id.String()).
			Return(expectedUser, nil)

		mocks.storageRepo.EXPECT().
			GetSignedURL(avatarURL).
			Return(signedURL, nil)

		user, err := svc.GetUserByID(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, signedURL, *user.AvatarURL)
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", id.String()).
			Return(nil, errors.New("user not found: blablabla"))

		user, err := svc.GetUserByID(ctx, id)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - repository error", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", id.String()).
			Return(nil, errors.New("db error"))

		user, err := svc.GetUserByID(ctx, id)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - storage signed URL error", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		avatarURL := "storage/avatar.jpg"
		expectedUser := &entity.User{
			ID:        id,
			Name:      "Test User",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			AvatarURL: &avatarURL,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", id.String()).
			Return(expectedUser, nil)

		mocks.storageRepo.EXPECT().
			GetSignedURL(avatarURL).
			Return("", errors.New("storage error"))

		user, err := svc.GetUserByID(ctx, id)
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_UserService_UpdatePassword(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	newPassword := "newPassword123"
	hashedPassword := "hashed_new_password"

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		oldPasswordHash := "old_password_hash"
		existingUser := &entity.User{
			ID:           uuid.New(),
			Name:         "Test User",
			Email:        email,
			PasswordHash: oldPasswordHash,
			Role:         enum.RoleStudent,
		}

		// Expect to get user by email
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(existingUser, nil)

		// Expect password hashing
		mocks.bcrypt.EXPECT().
			Hash(newPassword).
			Return(hashedPassword, nil)

		// Expect user update
		updatedUser := *existingUser
		updatedUser.PasswordHash = hashedPassword
		mocks.userRepo.EXPECT().
			UpdateUser(ctx, &updatedUser).
			Return(nil)

		err := svc.UpdatePassword(ctx, email, newPassword)
		assert.NoError(t, err)
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		// Expect to get user by email - returns not found
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(nil, errors.New("user not found: blablabla"))

		err := svc.UpdatePassword(ctx, email, newPassword)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - password hashing fails", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		oldPasswordHash := "old_password_hash"
		existingUser := &entity.User{
			ID:           uuid.New(),
			Name:         "Test User",
			Email:        email,
			PasswordHash: oldPasswordHash,
			Role:         enum.RoleStudent,
		}

		// Expect to get user by email
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(existingUser, nil)

		// Expect password hashing to fail
		mocks.bcrypt.EXPECT().
			Hash(newPassword).
			Return("", errors.New("hashing error"))

		err := svc.UpdatePassword(ctx, email, newPassword)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - update user fails", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		oldPasswordHash := "old_password_hash"
		existingUser := &entity.User{
			ID:           uuid.New(),
			Name:         "Test User",
			Email:        email,
			PasswordHash: oldPasswordHash,
			Role:         enum.RoleStudent,
		}

		// Expect to get user by email
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "email", email).
			Return(existingUser, nil)

		// Expect password hashing
		mocks.bcrypt.EXPECT().
			Hash(newPassword).
			Return(hashedPassword, nil)

		// Expect user update to fail
		updatedUser := *existingUser
		updatedUser.PasswordHash = hashedPassword
		mocks.userRepo.EXPECT().
			UpdateUser(ctx, &updatedUser).
			Return(errors.New("db error"))

		err := svc.UpdatePassword(ctx, email, newPassword)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_UserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	name := "Updated Name"
	bio := "Updated Bio"

	t.Run("success - update all fields", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Bio:       nil,
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := dto.UpdateUserRequest{
			Name: &name,
			Bio:  &bio,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Expect user update with new values
		expectedUpdatedUser := *existingUser
		expectedUpdatedUser.Name = name
		expectedUpdatedUser.Bio = &bio
		mocks.userRepo.EXPECT().
			UpdateUser(ctx, &expectedUpdatedUser).
			Return(nil)

		err := svc.UpdateUser(ctx, userID, req)
		assert.NoError(t, err)
	})

	t.Run("success - update partial fields", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Bio:       nil,
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := dto.UpdateUserRequest{
			Name: &name,
			Bio:  nil,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Expect user update with only name updated
		expectedUpdatedUser := *existingUser
		expectedUpdatedUser.Name = name
		mocks.userRepo.EXPECT().
			UpdateUser(ctx, &expectedUpdatedUser).
			Return(nil)

		err := svc.UpdateUser(ctx, userID, req)
		assert.NoError(t, err)
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := dto.UpdateUserRequest{
			Name: &name,
			Bio:  &bio,
		}

		// Expect to get user by ID - returns not found
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(nil, errors.New("user not found: blablabla"))

		err := svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - repository error during get", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		req := dto.UpdateUserRequest{
			Name: &name,
			Bio:  &bio,
		}

		// Expect to get user by ID - returns error
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(nil, errors.New("db error"))

		err := svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - repository error during update", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Bio:       nil,
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		req := dto.UpdateUserRequest{
			Name: &name,
			Bio:  &bio,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Expect user update to fail
		expectedUpdatedUser := *existingUser
		expectedUpdatedUser.Name = name
		expectedUpdatedUser.Bio = &bio
		mocks.userRepo.EXPECT().
			UpdateUser(ctx, &expectedUpdatedUser).
			Return(errors.New("db error"))

		err := svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("success - update avatar", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Bio:       nil,
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create a multipart form file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("avatar", "test.jpg")
		assert.NoError(t, err)

		// Write some test image content
		imageContent := []byte("mock image content")
		_, err = part.Write(imageContent)
		assert.NoError(t, err)
		writer.Close()

		// Read the form back to get the FileHeader
		reader := multipart.NewReader(body, writer.Boundary())
		form, err := reader.ReadForm(2 * fileutil.MegaByte)
		assert.NoError(t, err)

		avatar := form.File["avatar"][0]
		req := dto.UpdateUserRequest{
			Avatar: avatar,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Mock MIME type checking
		mocks.fileUtil.EXPECT().
			CheckMIMEFileType(mock.Anything, fileutil.ImageContentTypes).
			Return(true, "image/jpeg", nil)

		// Mock storage upload
		expectedAvatarURL := "https://storage.example.com/users/avatar/" + userID.String()
		mocks.storageRepo.EXPECT().
			Upload(ctx, mock.Anything, fmt.Sprintf("users/avatar/%s", userID.String())).
			Return(expectedAvatarURL, nil)

		// Expect user update with new avatar URL
		expectedUpdatedUser := *existingUser
		expectedUpdatedUser.AvatarURL = &expectedAvatarURL
		mocks.userRepo.EXPECT().
			UpdateUser(ctx, &expectedUpdatedUser).
			Return(nil)

		err = svc.UpdateUser(ctx, userID, req)
		assert.NoError(t, err)
	})

	t.Run("error - avatar file too large", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create a multipart form file that's too large
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("avatar", "large.jpg")
		assert.NoError(t, err)

		// Write content that exceeds the limit
		largeContent := bytes.Repeat([]byte("a"), 3*int(fileutil.MegaByte))
		_, err = part.Write(largeContent)
		assert.NoError(t, err)
		writer.Close()

		// Read the form back to get the FileHeader
		reader := multipart.NewReader(body, writer.Boundary())
		form, err := reader.ReadForm(4 * fileutil.MegaByte)
		assert.NoError(t, err)

		avatar := form.File["avatar"][0]
		req := dto.UpdateUserRequest{
			Avatar: avatar,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		err = svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrFileTooLarge)
	})

	t.Run("error - invalid file format", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create a multipart form file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("avatar", "test.txt")
		assert.NoError(t, err)

		_, err = part.Write([]byte("text content"))
		assert.NoError(t, err)
		writer.Close()

		// Read the form back to get the FileHeader
		reader := multipart.NewReader(body, writer.Boundary())
		form, err := reader.ReadForm(2 * fileutil.MegaByte)
		assert.NoError(t, err)

		avatar := form.File["avatar"][0]
		req := dto.UpdateUserRequest{
			Avatar: avatar,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Mock MIME type checking - returns invalid format
		mocks.fileUtil.EXPECT().
			CheckMIMEFileType(mock.Anything, fileutil.ImageContentTypes).
			Return(false, "text/plain", nil)

		err = svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidFileFormat)
	})

	t.Run("error - storage upload failure", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create a multipart form file
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("avatar", "test.jpg")
		assert.NoError(t, err)

		_, err = part.Write([]byte("image content"))
		assert.NoError(t, err)
		writer.Close()

		// Read the form back to get the FileHeader
		reader := multipart.NewReader(body, writer.Boundary())
		form, err := reader.ReadForm(2 * fileutil.MegaByte)
		assert.NoError(t, err)

		avatar := form.File["avatar"][0]
		req := dto.UpdateUserRequest{
			Avatar: avatar,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Mock MIME type checking
		mocks.fileUtil.EXPECT().
			CheckMIMEFileType(mock.Anything, fileutil.ImageContentTypes).
			Return(true, "image/jpeg", nil)

		// Mock storage upload failure
		mocks.storageRepo.EXPECT().
			Upload(ctx, mock.Anything, fmt.Sprintf("users/avatar/%s", userID.String())).
			Return("", errors.New("storage error"))

		err = svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - failed to open avatar file", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		avatar := &multipart.FileHeader{
			Filename: "test.jpg",
			Size:     100,
		}

		req := dto.UpdateUserRequest{
			Avatar: avatar,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		err := svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - MIME type check failure", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		existingUser := &entity.User{
			ID:        userID,
			Name:      "Original Name",
			Email:     "test@example.com",
			Role:      enum.RoleStudent,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create a multipart form file that will cause MIME check to fail
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("avatar", "test.jpg")
		require.NoError(t, err)

		_, err = part.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF}) // Invalid bytes
		require.NoError(t, err)
		writer.Close()

		// Read the form back to get the FileHeader
		reader := multipart.NewReader(body, writer.Boundary())
		form, err := reader.ReadForm(2 * fileutil.MegaByte)
		require.NoError(t, err)

		avatar := form.File["avatar"][0]
		req := dto.UpdateUserRequest{
			Avatar: avatar,
		}

		// Expect to get user by ID
		mocks.userRepo.EXPECT().
			GetUserByField(ctx, "id", userID.String()).
			Return(existingUser, nil)

		// Mock MIME type checking to return error
		mocks.fileUtil.EXPECT().
			CheckMIMEFileType(mock.Anything, fileutil.ImageContentTypes).
			Return(false, "", errors.New("failed to check MIME type"))

		err = svc.UpdateUser(ctx, userID, req)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_UserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		// Expect user deletion
		mocks.userRepo.EXPECT().
			DeleteUser(ctx, userID).
			Return(nil)

		err := svc.DeleteUser(ctx, userID)
		assert.NoError(t, err)
	})

	t.Run("success - with requester ID in context", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		// Create context with requester ID
		requesterID := uuid.New()
		ctxWithUser := context.WithValue(ctx, ctxkey.UserID, requesterID)

		// Expect user deletion
		mocks.userRepo.EXPECT().
			DeleteUser(ctxWithUser, userID).
			Return(nil)

		err := svc.DeleteUser(ctxWithUser, userID)
		assert.NoError(t, err)
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		// Expect deletion to return not found error
		mocks.userRepo.EXPECT().
			DeleteUser(ctx, userID).
			Return(errors.New("user not found"))

		err := svc.DeleteUser(ctx, userID)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - repository error", func(t *testing.T) {
		svc, mocks := setupUserServiceTest(t)

		// Expect deletion to return generic error
		mocks.userRepo.EXPECT().
			DeleteUser(ctx, userID).
			Return(errors.New("db error"))

		err := svc.DeleteUser(ctx, userID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}
