package dto

import (
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/pkg/fileutil"
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
	u.CreatedAt = &user.CreatedAt
	u.UpdatedAt = &user.UpdatedAt

	if user.ID != uuid.Nil && user.AvatarURL != nil {
		signedURL := fileutil.GetSignedURL(fmt.Sprintf("users/avatar/%s", user.ID.String()))
		u.AvatarURL = &signedURL
	}

	return u
}

func (u *UserResponse) PopulateMinimalFromEntity(user *entity.User) *UserResponse {
	u.ID = user.ID
	u.Name = user.Name
	u.Role = user.Role
	u.Bio = user.Bio

	if user.ID != uuid.Nil && user.AvatarURL != nil {
		signedURL := fileutil.GetSignedURL(fmt.Sprintf("users/avatar/%s", user.ID.String()))
		u.AvatarURL = &signedURL
	}

	return u
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string `json:"-"`
	Role     enum.UserRole
}

type UpdateUserRequest struct {
	Name   *string               `form:"name" validate:"omitempty,min=3,max=100,ascii"`
	Bio    *string               `form:"bio"  validate:"omitempty,max=500"`
	Avatar *multipart.FileHeader `form:"avatar"`
}
