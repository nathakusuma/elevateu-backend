package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/app/auth/service"
	appmocks "github.com/nathakusuma/elevateu-backend/test/unit/mocks/app"
	pkgmocks "github.com/nathakusuma/elevateu-backend/test/unit/mocks/pkg"
	_ "github.com/nathakusuma/elevateu-backend/test/unit/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type authServiceMocks struct {
	authRepo *appmocks.MockIAuthRepository
	userSvc  *appmocks.MockIUserService
	bcrypt   *pkgmocks.MockIBcrypt
	jwt      *pkgmocks.MockIJwt
	mailer   *pkgmocks.MockIMailer
	radgen   *pkgmocks.MockIRandGen
	uuid     *pkgmocks.MockIUUID
}

func setupAuthServiceMocks(t *testing.T) (contract.IAuthService, *authServiceMocks) {
	mocks := &authServiceMocks{
		authRepo: appmocks.NewMockIAuthRepository(t),
		userSvc:  appmocks.NewMockIUserService(t),
		bcrypt:   pkgmocks.NewMockIBcrypt(t),
		jwt:      pkgmocks.NewMockIJwt(t),
		mailer:   pkgmocks.NewMockIMailer(t),
		radgen:   pkgmocks.NewMockIRandGen(t),
		uuid:     pkgmocks.NewMockIUUID(t),
	}

	svc := service.NewAuthService(mocks.authRepo, mocks.userSvc, mocks.bcrypt,
		mocks.jwt, mocks.mailer, mocks.radgen, mocks.uuid)

	return svc, mocks
}

func Test_AuthService_RequestOTPRegister(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)
		emailSent := make(chan struct{}, 1)

		// Expect user not found (which is good for registration)
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errorpkg.ErrNotFound)

		// Expect OTP generation success
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(123456, nil)

		// Expect OTP to be set
		mocks.authRepo.EXPECT().
			SetRegisterOTP(ctx, email, mock.AnythingOfType("string")).
			Return(nil)

		// Mock email sending with channel notification
		mocks.mailer.EXPECT().
			Send(
				email,
				"[ElevateU] Verify Your Account",
				"otp_register.html",
				mock.AnythingOfType("map[string]interface {}"),
			).RunAndReturn(func(_, _, _ string, _ map[string]interface{}) error {
			emailSent <- struct{}{}
			return nil
		})

		err := svc.RequestRegisterOTP(ctx, email)
		assert.NoError(t, err)

		// Wait for email sending goroutine to complete
		<-emailSent
	})

	t.Run("error - otp generation fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect user not found (which is good for registration)
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errorpkg.ErrNotFound)

		// Expect OTP generation error
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(0, errors.New("otp generation error"))

		err := svc.RequestRegisterOTP(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - email sending fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)
		emailSent := make(chan struct{}, 1)

		// Expect user not found (which is good for registration)
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errorpkg.ErrNotFound)

		// Expect OTP generation success
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(123456, nil)

		// Expect OTP to be set
		mocks.authRepo.EXPECT().
			SetRegisterOTP(ctx, email, mock.AnythingOfType("string")).
			Return(nil)

		// Mock email sending to fail
		mocks.mailer.EXPECT().
			Send(
				email,
				"[ElevateU] Verify Your Account",
				"otp_register.html",
				mock.AnythingOfType("map[string]interface {}"),
			).RunAndReturn(func(_, _, _ string, _ map[string]interface{}) error {
			emailSent <- struct{}{}
			return errors.New("email sending error")
		})

		err := svc.RequestRegisterOTP(ctx, email)
		assert.NoError(t, err)

		// Wait for email sending goroutine to complete
		<-emailSent

		// It should not return an error because the email sending is done in a goroutine
	})

	t.Run("error - email already registered", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Return existing user
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(&entity.User{ID: uuid.New()}, nil)

		err := svc.RequestRegisterOTP(ctx, email)
		assert.ErrorIs(t, err, errorpkg.ErrEmailAlreadyRegistered)
	})

	t.Run("error - get user unexpected error", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errors.New("unexpected error"))

		err := svc.RequestRegisterOTP(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - set OTP fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errorpkg.ErrNotFound)

		// Expect OTP generation success
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(123456, nil)

		mocks.authRepo.EXPECT().
			SetRegisterOTP(ctx, email, mock.AnythingOfType("string")).
			Return(errors.New("redis error"))

		err := svc.RequestRegisterOTP(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_AuthService_CheckRegisterOTP(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"
	otp := "123456"

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, email).
			Return(otp, nil)

		err := svc.CheckRegisterOTP(ctx, email, otp)
		assert.NoError(t, err)
	})

	t.Run("error - OTP not found", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, email).
			Return("", errors.New("otp not found: blablabla"))

		err := svc.CheckRegisterOTP(ctx, email, otp)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidOTP)
	})

	t.Run("error - get OTP fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, email).
			Return("", errors.New("redis error"))

		err := svc.CheckRegisterOTP(ctx, email, otp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - invalid OTP", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, email).
			Return("654321", nil)

		err := svc.CheckRegisterOTP(ctx, email, otp)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidOTP)
	})
}

func Test_AuthService_Login(t *testing.T) {
	ctx := context.Background()
	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		userID := uuid.New()
		user := &entity.User{
			Email: req.Email,
		}

		mockLoginExpectations(ctx, mocks, req.Email, req.Password, userID)

		resp, err := svc.Login(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.Equal(t, user.Email, resp.User.Email)
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(nil, errorpkg.ErrNotFound)

		resp, err := svc.Login(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - get user unexpected error", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(nil, errors.New("db error"))

		resp, err := svc.Login(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - invalid credentials", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		passwordHash := "hashed_password"
		user := &entity.User{
			Email:        req.Email,
			PasswordHash: passwordHash,
		}

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(user, nil)

		mocks.bcrypt.EXPECT().
			Compare(req.Password, passwordHash).
			Return(false)

		resp, err := svc.Login(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrCredentialsNotMatch)
	})

	t.Run("error - jwt creation fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		passwordHash := "hashed_password"
		user := &entity.User{
			ID:           uuid.New(),
			Email:        req.Email,
			PasswordHash: passwordHash,
			Role:         enum.UserRoleStudent,
		}

		// Setup expectations
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(user, nil)

		mocks.bcrypt.EXPECT().
			Compare(req.Password, passwordHash).
			Return(true)

		// JWT creation will fail
		mocks.jwt.EXPECT().
			Create(user.ID, user.Role).
			Return("", errors.New("jwt error"))

		// Expect CreateAuthSession to be called but we don't care about the result
		// since the JWT error should be returned first
		mocks.authRepo.EXPECT().
			CreateAuthSession(ctx, mock.MatchedBy(func(authSession *entity.AuthSession) bool {
				return authSession.UserID == user.ID &&
					len(authSession.Token) == 32 &&
					!authSession.ExpiresAt.IsZero()
			})).
			Return(nil).
			Maybe() // This may or may not be called depending on goroutine scheduling

		resp, err := svc.Login(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - refresh token generation fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		passwordHash := "hashed_password"
		user := &entity.User{
			ID:           uuid.New(),
			Email:        req.Email,
			PasswordHash: passwordHash,
			Role:         enum.UserRoleStudent,
		}

		// Setup expectations
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(user, nil)

		mocks.bcrypt.EXPECT().
			Compare(req.Password, passwordHash).
			Return(true)

		// JWT creation succeeds
		mocks.jwt.EXPECT().
			Create(user.ID, user.Role).
			Return("access_token", nil)

		// Refresh token generation fails
		mocks.radgen.EXPECT().
			RandomString(32).
			Return("", errors.New("randgen error"))

		resp, err := svc.Login(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - create auth session fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		passwordHash := "hashed_password"
		user := &entity.User{
			ID:           uuid.New(),
			Email:        req.Email,
			PasswordHash: passwordHash,
			Role:         enum.UserRoleStudent,
		}

		// Setup expectations
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(user, nil)

		mocks.bcrypt.EXPECT().
			Compare(req.Password, passwordHash).
			Return(true)

		// JWT creation succeeds
		mocks.jwt.EXPECT().
			Create(user.ID, user.Role).
			Return("access_token", nil)

		// Refresh token generation succeeds
		mocks.radgen.EXPECT().
			RandomString(32).
			Return("12345678901234567890123456789012", nil)

		// AuthSession creation fails
		mocks.authRepo.EXPECT().
			CreateAuthSession(ctx, mock.MatchedBy(func(authSession *entity.AuthSession) bool {
				return authSession.UserID == user.ID &&
					len(authSession.Token) == 32 &&
					!authSession.ExpiresAt.IsZero()
			})).
			Return(errors.New("db error"))

		resp, err := svc.Login(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_AuthService_Register(t *testing.T) {
	ctx := context.Background()
	req := dto.RegisterRequest{
		Email:    "test@example.com",
		OTP:      "123456",
		Name:     "Test User",
		Password: "password123",
	}

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Setup expectations
		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeleteRegisterOTP(ctx, req.Email).
			Return(nil)

		userID := uuid.New()
		mocks.userSvc.EXPECT().
			CreateUser(ctx, &dto.CreateUserRequest{
				Name:     req.Name,
				Email:    req.Email,
				Password: req.Password,
				Role:     enum.UserRoleStudent,
			}).Return(userID, nil)

		// Mock login expectations
		mockLoginExpectations(ctx, mocks, req.Email, req.Password, userID)

		resp, err := svc.Register(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotNil(t, resp.User)
	})

	t.Run("error - no OTP found", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, req.Email).
			Return("", errors.New("otp not found: blablabla"))

		resp, err := svc.Register(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidOTP)
	})

	t.Run("error - invalid OTP", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, req.Email).
			Return("different-otp", nil)

		resp, err := svc.Register(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidOTP)
	})

	t.Run("error - get OTP fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, req.Email).
			Return("", errors.New("redis error"))

		resp, err := svc.Register(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - delete OTP fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeleteRegisterOTP(ctx, req.Email).
			Return(errors.New("redis error"))

		resp, err := svc.Register(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - create user fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetRegisterOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeleteRegisterOTP(ctx, req.Email).
			Return(nil)

		mocks.userSvc.EXPECT().
			CreateUser(ctx, mock.Anything).
			Return(uuid.UUID{}, errors.New("db error"))

		resp, err := svc.Register(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

// Helper function to set up common login expectations
func mockLoginExpectations(ctx context.Context, mocks *authServiceMocks, email, password string, userID uuid.UUID) {
	passwordHash := "hashed_password"
	user := &entity.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         enum.UserRoleStudent,
		Name:         "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mocks.userSvc.EXPECT().
		GetUserByEmail(ctx, email).
		Return(user, nil)

	mocks.bcrypt.EXPECT().
		Compare(password, passwordHash).
		Return(true)

	mocks.jwt.EXPECT().
		Create(user.ID, user.Role).
		Return("access_token", nil)

	mocks.radgen.EXPECT().
		RandomString(32).
		Return("12345678901234567890123456789012", nil)

	mocks.authRepo.EXPECT().
		CreateAuthSession(ctx, mock.MatchedBy(func(authSession *entity.AuthSession) bool {
			return authSession.UserID == user.ID &&
				len(authSession.Token) == 32 && // Check refresh token length
				!authSession.ExpiresAt.IsZero() // Check expiration is set
		})).
		Return(nil)
}

func Test_AuthService_Refresh(t *testing.T) {
	ctx := context.Background()
	refreshToken := "test-refresh-token"
	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"
	userID := uuid.New()

	mockUser := entity.User{
		ID:        userID,
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      enum.UserRoleStudent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockAuthSession := entity.AuthSession{
		Token:     refreshToken,
		UserID:    userID,
		User:      mockUser,
		ExpiresAt: time.Now().Add(time.Hour), // Valid for 1 hour
	}

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect to get valid auth session
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(&mockAuthSession, nil)

		// Expect token generation
		mocks.jwt.EXPECT().
			Create(userID, enum.UserRoleStudent).
			Return(newAccessToken, nil)

		// Expect random string generation for refresh token
		mocks.radgen.EXPECT().
			RandomString(32).
			Return(newRefreshToken, nil)

		// Expect new auth session creation
		mocks.authRepo.EXPECT().
			CreateAuthSession(ctx, mock.AnythingOfType("*entity.AuthSession")).
			Return(nil)

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert response
		assert.NoError(t, err)
		assert.Equal(t, newAccessToken, resp.AccessToken)
		assert.Equal(t, newRefreshToken, resp.RefreshToken)
		assert.Equal(t, mockUser.ID, resp.User.ID)
		assert.Equal(t, mockUser.Email, resp.User.Email)
		assert.Equal(t, mockUser.Role, resp.User.Role)
	})

	t.Run("error - invalid refresh token", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect auth session not found
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(nil, errors.New("auth session not found"))

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidRefreshToken)
		assert.Empty(t, resp)
	})

	t.Run("error - expired refresh token", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		expiredSession := mockAuthSession
		expiredSession.ExpiresAt = time.Now().Add(-time.Hour) // Expired 1 hour ago

		// Expect to get expired auth session
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(&expiredSession, nil)

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidRefreshToken)
		assert.Empty(t, resp)
	})

	t.Run("error - failed to get auth session", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect unexpected error when getting auth session
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(nil, errors.New("unexpected error"))

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
		assert.Empty(t, resp)
	})

	t.Run("error - failed to generate access token", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect to get valid auth session
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(&mockAuthSession, nil)

		// Expect error in token generation
		mocks.jwt.EXPECT().
			Create(userID, enum.UserRoleStudent).
			Return("", errors.New("token generation failed"))

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
		assert.Empty(t, resp)
	})

	t.Run("error - failed to generate refresh token", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect to get valid auth session
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(&mockAuthSession, nil)

		// Expect successful access token generation
		mocks.jwt.EXPECT().
			Create(userID, enum.UserRoleStudent).
			Return(newAccessToken, nil)

		// Expect error in refresh token generation
		mocks.radgen.EXPECT().
			RandomString(32).
			Return("", errors.New("random generation failed"))

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
		assert.Empty(t, resp)
	})

	t.Run("error - failed to create auth session", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Expect to get valid auth session
		mocks.authRepo.EXPECT().
			GetAuthSessionByToken(ctx, refreshToken).
			Return(&mockAuthSession, nil)

		// Expect successful access token generation
		mocks.jwt.EXPECT().
			Create(userID, enum.UserRoleStudent).
			Return(newAccessToken, nil)

		// Expect successful refresh token generation
		mocks.radgen.EXPECT().
			RandomString(32).
			Return(newRefreshToken, nil)

		// Expect error in creating auth session
		mocks.authRepo.EXPECT().
			CreateAuthSession(ctx, mock.AnythingOfType("*entity.AuthSession")).
			Return(errors.New("failed to create auth session"))

		// Call the service
		resp, err := svc.Refresh(ctx, refreshToken)

		// Assert error
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
		assert.Empty(t, resp)
	})
}

func Test_AuthService_Logout(t *testing.T) {
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), ctxkey.UserID, userID)

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			DeleteAuthSession(ctx, userID).
			Return(nil)

		err := svc.Logout(ctx, userID)
		assert.NoError(t, err)
	})

	t.Run("error - delete auth session fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			DeleteAuthSession(ctx, userID).
			Return(errors.New("db error"))

		err := svc.Logout(ctx, userID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}

func Test_AuthService_RequestPasswordResetOTP(t *testing.T) {
	ctx := context.Background()
	email := "test@example.com"

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)
		emailSent := make(chan struct{}, 1)

		// Expect user to be found
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(&entity.User{
				ID:           uuid.New(),
				PasswordHash: "hashed_password",
			}, nil)

		// Expect OTP generation success
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(123456, nil)

		// Expect OTP to be set
		mocks.authRepo.EXPECT().
			SetPasswordResetOTP(ctx, email, mock.AnythingOfType("string")).
			Return(nil)

		// Mock email sending with channel notification
		mocks.mailer.EXPECT().
			Send(
				email,
				"[ElevateU] Reset Password",
				"otp_reset_password.html",
				mock.AnythingOfType("map[string]interface {}"),
			).RunAndReturn(func(_, _, _ string, _ map[string]interface{}) error {
			emailSent <- struct{}{}
			return nil
		})

		err := svc.RequestPasswordResetOTP(ctx, email)
		assert.NoError(t, err)

		// Wait for email sending goroutine to complete
		<-emailSent
	})

	t.Run("error - user not found", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errorpkg.ErrNotFound)

		err := svc.RequestPasswordResetOTP(ctx, email)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - get user unexpected error", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(nil, errors.New("unexpected error"))

		err := svc.RequestPasswordResetOTP(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - OTP generation fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(&entity.User{
				ID:           uuid.New(),
				PasswordHash: "hashed_password",
			}, nil)

		// Expect OTP generation error
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(0, errors.New("otp generation error"))

		err := svc.RequestPasswordResetOTP(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - set OTP fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(&entity.User{
				ID:           uuid.New(),
				PasswordHash: "hashed_password",
			}, nil)

		// Expect OTP generation success
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(123456, nil)

		mocks.authRepo.EXPECT().
			SetPasswordResetOTP(ctx, email, mock.AnythingOfType("string")).
			Return(errors.New("redis error"))

		err := svc.RequestPasswordResetOTP(ctx, email)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - email sending fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)
		emailSent := make(chan struct{}, 1)

		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, email).
			Return(&entity.User{
				ID:           uuid.New(),
				PasswordHash: "hashed_password",
			}, nil)

		// Expect OTP generation success
		mocks.radgen.EXPECT().
			RandomNumber(6).
			Return(123456, nil)

		mocks.authRepo.EXPECT().
			SetPasswordResetOTP(ctx, email, mock.AnythingOfType("string")).
			Return(nil)

		mocks.mailer.EXPECT().
			Send(
				email,
				"[ElevateU] Reset Password",
				"otp_reset_password.html",
				mock.AnythingOfType("map[string]interface {}"),
			).RunAndReturn(func(_, _, _ string, _ map[string]interface{}) error {
			emailSent <- struct{}{}
			return errors.New("email sending error")
		})

		err := svc.RequestPasswordResetOTP(ctx, email)
		assert.NoError(t, err) // Should not return error as email is sent in goroutine

		// Wait for email sending goroutine to complete
		<-emailSent
	})
}

func Test_AuthService_ResetPassword(t *testing.T) {
	ctx := context.Background()
	req := dto.ResetPasswordRequest{
		Email:       "test@example.com",
		OTP:         "123456",
		NewPassword: "newpassword123",
	}

	t.Run("success", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Setup expectations for password reset
		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeletePasswordResetOTP(ctx, req.Email).
			Return(nil)

		mocks.userSvc.EXPECT().
			UpdatePassword(ctx, req.Email, req.NewPassword).
			Return(nil)

		// Setup expectations for subsequent login
		mockLoginExpectations(ctx, mocks, req.Email, req.NewPassword, uuid.New())

		resp, err := svc.ResetPassword(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)
		assert.NotNil(t, resp.User)
	})

	t.Run("error - OTP not found", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return("", errors.New("otp not found: blablabla"))

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidOTP)
	})

	t.Run("error - get OTP unexpected error", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return("", errors.New("redis error"))

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - invalid OTP", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return("different-otp", nil)

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrInvalidOTP)
	})

	t.Run("error - delete OTP fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeletePasswordResetOTP(ctx, req.Email).
			Return(errors.New("redis error"))

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - update password fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeletePasswordResetOTP(ctx, req.Email).
			Return(nil)

		mocks.userSvc.EXPECT().
			UpdatePassword(ctx, req.Email, req.NewPassword).
			Return(errors.New("db error"))

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})

	t.Run("error - user not found during update", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeletePasswordResetOTP(ctx, req.Email).
			Return(nil)

		mocks.userSvc.EXPECT().
			UpdatePassword(ctx, req.Email, req.NewPassword).
			Return(errorpkg.ErrNotFound)

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.ErrorIs(t, err, errorpkg.ErrNotFound)
	})

	t.Run("error - login after reset fails", func(t *testing.T) {
		svc, mocks := setupAuthServiceMocks(t)

		// Success up to password update
		mocks.authRepo.EXPECT().
			GetPasswordResetOTP(ctx, req.Email).
			Return(req.OTP, nil)

		mocks.authRepo.EXPECT().
			DeletePasswordResetOTP(ctx, req.Email).
			Return(nil)

		mocks.userSvc.EXPECT().
			UpdatePassword(ctx, req.Email, req.NewPassword).
			Return(nil)

		// Fail during login
		mocks.userSvc.EXPECT().
			GetUserByEmail(ctx, req.Email).
			Return(nil, errors.New("db error"))

		resp, err := svc.ResetPassword(ctx, req)
		assert.Empty(t, resp)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errorpkg.ErrInternalServer)
	})
}
