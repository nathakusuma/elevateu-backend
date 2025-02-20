package contract

import (
	"context"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type IAuthRepository interface {
	SetRegisterOTP(ctx context.Context, email, otp string) error
	GetRegisterOTP(ctx context.Context, email string) (string, error)
	DeleteRegisterOTP(ctx context.Context, email string) error

	CreateAuthSession(ctx context.Context, authSession *entity.AuthSession) error
	GetAuthSessionByToken(ctx context.Context, token string) (*entity.AuthSession, error)
	DeleteAuthSession(ctx context.Context, userID uuid.UUID) error

	SetPasswordResetOTP(ctx context.Context, email, otp string) error
	GetPasswordResetOTP(ctx context.Context, email string) (string, error)
	DeletePasswordResetOTP(ctx context.Context, email string) error
}

type IAuthService interface {
	RequestRegisterOTP(ctx context.Context, email string) error
	CheckRegisterOTP(ctx context.Context, email, otp string) error
	Register(ctx context.Context, req dto.RegisterRequest) (dto.LoginResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)

	Refresh(ctx context.Context, refreshToken string) (dto.LoginResponse, error)
	Logout(ctx context.Context, userID uuid.UUID) error

	RequestPasswordResetOTP(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) (dto.LoginResponse, error)
}
