package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type UserResponse struct {
	ID        uuid.UUID     `json:"id,omitempty"`
	Name      string        `json:"name,omitempty"`
	Email     string        `json:"email,omitempty"`
	Role      enum.UserRole `json:"role,omitempty"`
	Bio       *string       `json:"bio,omitempty"`
	AvatarURL *string       `json:"avatar_url,omitempty"`
	CreatedAt *time.Time    `json:"created_at,omitempty"`
	UpdatedAt *time.Time    `json:"updated_at,omitempty"`
}

func (u *UserResponse) PopulateFromEntity(user *entity.User) *UserResponse {
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
	u.Role = user.Role
	u.Bio = user.Bio
	u.AvatarURL = user.AvatarURL
	u.CreatedAt = &user.CreatedAt
	u.UpdatedAt = &user.UpdatedAt

	return u
}

func (u *UserResponse) PopulateMinimalFromEntity(user *entity.User) *UserResponse {
	u.ID = user.ID
	u.Name = user.Name
	u.Role = user.Role
	u.Bio = user.Bio
	u.AvatarURL = user.AvatarURL

	return u
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string `json:"-"`
	Role     enum.UserRole
}

type UpdateUserRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=100,ascii"`
	Bio  *string `json:"bio"  validate:"omitempty,max=500"`
}
