package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/cache"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
	"github.com/nathakusuma/elevateu-backend/pkg/bcrypt"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
	"github.com/nathakusuma/elevateu-backend/pkg/jwt"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/mail"
	"github.com/nathakusuma/elevateu-backend/pkg/randgen"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type authService struct {
	repo     contract.IAuthRepository
	userSvc  contract.IUserService
	bcrypt   bcrypt.IBcrypt
	cache    cache.ICache
	fileUtil fileutil.IFileUtil
	jwt      jwt.IJwt
	mailer   mail.IMailer
	randgen  randgen.IRandGen
	uuid     uuidpkg.IUUID
}

func NewAuthService(
	authRepo contract.IAuthRepository,
	userSvc contract.IUserService,
	bcrypt bcrypt.IBcrypt,
	cache cache.ICache,
	fileUtil fileutil.IFileUtil,
	jwt jwt.IJwt,
	mailer mail.IMailer,
	randgen randgen.IRandGen,
	uuid uuidpkg.IUUID,
) contract.IAuthService {
	return &authService{
		repo:     authRepo,
		userSvc:  userSvc,
		bcrypt:   bcrypt,
		cache:    cache,
		fileUtil: fileUtil,
		jwt:      jwt,
		mailer:   mailer,
		randgen:  randgen,
		uuid:     uuid,
	}
}

func (s *authService) RequestRegisterOTP(ctx context.Context, email string) error {
	// check if email is already registered
	_, err := s.userSvc.GetUserByEmail(ctx, email)
	if err == nil {
		return errorpkg.ErrEmailAlreadyRegistered()
	}

	if !errors.Is(err, errorpkg.ErrNotFound()) {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "failed to get user by email")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// generate otp
	otpInt, err := s.randgen.RandomNumber(6)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "failed to generate otp")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	otp := strconv.Itoa(otpInt)

	// save otp
	err = s.cache.Set(ctx, "auth:"+email+":register_otp", otp, 10*time.Minute)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "failed to save otp")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// send otp to email
	go func() {
		err = s.mailer.Send(
			email,
			"[ElevateU] Verify Your Account",
			"otp_register.html",
			map[string]interface{}{
				"otp": otp,
			})
		if err != nil {
			log.Error(ctx, map[string]interface{}{
				"error": err,
			}, "failed to send email")
		}
	}()

	log.Info(ctx, map[string]interface{}{
		"user.email": email,
	}, "otp requested")

	return nil
}

func (s *authService) Register(ctx context.Context,
	req dto.RegisterRequest,
) (dto.LoginResponse, error) {
	var resp dto.LoginResponse

	// req without Password and OTP
	loggableReq := req
	loggableReq.Password = ""
	loggableReq.OTP = ""

	// get otp
	var savedOtp string
	err := s.cache.Get(ctx, "auth:"+req.Email+":register_otp", &savedOtp)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not found") {
			return resp, errorpkg.ErrInvalidOTP()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": loggableReq,
		}, "failed to get otp")
		return resp, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if savedOtp != req.OTP {
		return resp, errorpkg.ErrInvalidOTP()
	}

	// delete otp
	err = s.cache.Del(ctx, "auth:"+req.Email+":register_otp")
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": loggableReq,
		}, "failed to delete otp")
		return resp, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// Prepare user creation request
	createUserReq := &dto.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	// Add role-specific data
	if req.Role == enum.UserRoleStudent && req.Student != nil {
		createUserReq.Student = req.Student
	} else if req.Role == enum.UserRoleMentor && req.Mentor != nil {
		createUserReq.Mentor = req.Mentor
	}

	// save user
	_, err = s.userSvc.CreateUser(ctx, createUserReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": loggableReq,
		}, "failed to create user")
		return resp, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"user.email": req.Email,
		"user.role":  req.Role,
	}, "user registered")

	// login
	return s.Login(ctx, dto.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	// get user by email
	user, err := s.userSvc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errorpkg.ErrNotFound()) {
			return dto.LoginResponse{}, errorpkg.ErrNotFound().WithDetail("User not found. Please register first.")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": req.Email,
		}, "failed to get user by email")
		return dto.LoginResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// check password
	ok := s.bcrypt.Compare(req.Password, user.PasswordHash)
	if !ok {
		return dto.LoginResponse{}, errorpkg.ErrCredentialsNotMatch()
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	userResp := &dto.UserResponse{}
	if err = userResp.PopulateFromEntity(user, s.fileUtil.GetSignedURL); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
			"user":  user,
		}, "failed to populate user response")
		return dto.LoginResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"user.id":    user.ID,
		"user.email": user.Email,
	}, "user logged in")

	return dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResp,
	}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (dto.LoginResponse, error) {
	var resp dto.LoginResponse

	authSession, err := s.repo.GetAuthSessionByToken(ctx, refreshToken)
	if err != nil {
		if strings.HasPrefix(err.Error(), "auth session not found") {
			return resp, errorpkg.ErrInvalidRefreshToken()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "failed to get auth session by token")
		return resp, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if authSession.ExpiresAt.Before(time.Now()) {
		return resp, errorpkg.ErrInvalidRefreshToken()
	}

	// rotate refresh token
	accessToken, refreshToken, err := s.generateTokens(ctx, &authSession.User)
	if err != nil {
		return resp, err
	}

	userResp := &dto.UserResponse{}
	if err = userResp.PopulateFromEntity(&authSession.User, s.fileUtil.GetSignedURL); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "failed to populate user response")
		return resp, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	resp = dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResp,
	}

	log.Info(ctx, nil, "token refreshed")

	return resp, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.DeleteAuthSession(ctx, userID)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "failed to delete auth session")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, nil, "user logged out")

	return nil
}

func (s *authService) RequestPasswordResetOTP(ctx context.Context, email string) error {
	// check if email is registered
	_, err := s.userSvc.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, errorpkg.ErrNotFound()) {
			return errorpkg.ErrNotFound().WithDetail("User not found. Please register.")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "failed to get user by email")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// generate otp
	otpInt, err := s.randgen.RandomNumber(6)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "failed to generate otp")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	otp := strconv.Itoa(otpInt)

	// save otp
	err = s.cache.Set(ctx, "auth:"+email+":reset_password_otp", otp, 10*time.Minute)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": email,
		}, "failed to save otp")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// send otp to email
	go func() {
		err = s.mailer.Send(
			email,
			"[ElevateU] Reset Password",
			"otp_reset_password.html",
			map[string]interface{}{
				"otp": otp,
			})
		if err != nil {
			log.Error(ctx, map[string]interface{}{
				"error": err,
			}, "failed to send email")
		}
	}()

	log.Info(ctx, map[string]interface{}{
		"user.email": email,
	}, "otp requested")

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (dto.LoginResponse, error) {
	// get otp
	var savedOtp string
	err := s.cache.Get(ctx, "auth:"+req.Email+":reset_password_otp", &savedOtp)
	if err != nil {
		if strings.HasPrefix(err.Error(), "not found") {
			return dto.LoginResponse{}, errorpkg.ErrInvalidOTP()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": req.Email,
		}, "failed to get otp")
		return dto.LoginResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if savedOtp != req.OTP {
		return dto.LoginResponse{}, errorpkg.ErrInvalidOTP()
	}

	// delete otp
	err = s.cache.Del(ctx, "auth:"+req.Email+":reset_password_otp")
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": req.Email,
		}, "failed to delete otp")
		return dto.LoginResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// update user password
	if err = s.userSvc.UpdatePassword(ctx, req.Email, req.NewPassword); err != nil {
		if errors.Is(err, errorpkg.ErrNotFound()) {
			// Small chance, since we've already checked it on RequestPasswordResetOTP
			return dto.LoginResponse{}, err
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":      err,
			"user.email": req.Email,
		}, "failed to update user password")
		return dto.LoginResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"user.email": req.Email,
	}, "password reset")

	return s.Login(ctx, dto.LoginRequest{
		Email:    req.Email,
		Password: req.NewPassword,
	})
}

func (s *authService) generateTokens(ctx context.Context, user *entity.User) (string, string, error) {
	// Generate access token
	var isSubscribedBoost, isSubscribedChallenge bool
	if user.Role == enum.UserRoleStudent && user.Student != nil {
		isSubscribedBoost = user.Student.SubscribedBoostUntil.After(time.Now())
		isSubscribedChallenge = user.Student.SubscribedChallengeUntil.After(time.Now())
	}

	accessToken, err := s.jwt.Create(user.ID, user.Role, isSubscribedBoost, isSubscribedChallenge)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"user.id": user.ID,
		}, "Failed to generate access token")
		return "", "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// Generate refresh token
	refreshToken, err := s.randgen.RandomString(32)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"user.id": user.ID,
		}, "Failed to generate refresh token")
		return "", "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	// Create auth session
	if err = s.repo.CreateAuthSession(ctx, &entity.AuthSession{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(env.GetEnv().JwtRefreshExpireDuration),
	}); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"user.id": user.ID,
		}, "Failed to store auth session")
		return "", "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	return accessToken, refreshToken, nil
}
