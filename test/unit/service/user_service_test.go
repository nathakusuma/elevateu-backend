package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/app/user/service"
	appmocks "github.com/nathakusuma/elevateu-backend/test/unit/mocks/app"
	pkgmocks "github.com/nathakusuma/elevateu-backend/test/unit/mocks/pkg"
	_ "github.com/nathakusuma/elevateu-backend/test/unit/setup" // Initialize test environment
)

type userServiceMocks struct {
	userRepo *appmocks.MockIUserRepository
	uuid     *pkgmocks.MockIUUID
	bcrypt   *pkgmocks.MockIBcrypt
}

func setupUserServiceTest(t *testing.T) (contract.IUserService, *userServiceMocks) {
	mocks := &userServiceMocks{
		userRepo: appmocks.NewMockIUserRepository(t),
		uuid:     pkgmocks.NewMockIUUID(t),
		bcrypt:   pkgmocks.NewMockIBcrypt(t),
	}

	svc := service.NewUserService(mocks.userRepo, mocks.bcrypt, mocks.uuid)

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

	t.Run("success", func(t *testing.T) {
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
}

func Test_UserService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
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
