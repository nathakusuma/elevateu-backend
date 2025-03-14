package contract

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error)
	UpdateUser(ctx context.Context, req *dto.UserUpdate) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	AddPoint(ctx context.Context, txWrapper database.ITransaction, userID uuid.UUID, point int) error
	GetTopPoints(ctx context.Context, limit int) ([]*entity.User, error)
	GetMentors(ctx context.Context, pageReq dto.PaginationRequest) ([]*entity.User, dto.PaginationResponse, error)
}

type IUserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID, isMinimal bool) (*dto.UserResponse, error)
	UpdatePassword(ctx context.Context, email, newPassword string) error
	UpdateUser(ctx context.Context, id uuid.UUID, req dto.UpdateUserRequest) error
	UpdateUserAvatar(ctx context.Context, id uuid.UUID, avatar *multipart.FileHeader) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	DeleteUserAvatar(ctx context.Context, id uuid.UUID) error

	GetLeaderboard(ctx context.Context) ([]*dto.UserResponse, error)
	GetMentors(ctx context.Context, pageReq dto.PaginationRequest) ([]*dto.UserResponse, dto.PaginationResponse, error)
}
