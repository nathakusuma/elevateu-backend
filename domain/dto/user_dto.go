package dto

import (
	"mime/multipart"
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
	AvatarURL *string       `json:"avatar_url,omitempty"`
	CreatedAt *time.Time    `json:"created_at,omitempty"`
	UpdatedAt *time.Time    `json:"updated_at,omitempty"`
	Student   *StudentData  `json:"student,omitempty"`
	Mentor    *MentorData   `json:"mentor,omitempty"`
}

type StudentData struct {
	Instance string `json:"instance,omitempty"`
	Major    string `json:"major,omitempty"`
}

type MentorData struct {
	Specialization string  `json:"specialization,omitempty"`
	Experience     string  `json:"experience,omitempty"`
	Rating         float64 `json:"rating,omitempty"`
	RatingCount    int     `json:"rating_count,omitempty"`
	Price          int     `json:"price,omitempty"`
}

func (u *UserResponse) PopulateFromEntity(user *entity.User) *UserResponse {
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
	u.Role = user.Role
	u.AvatarURL = user.AvatarURL
	u.CreatedAt = &user.CreatedAt
	u.UpdatedAt = &user.UpdatedAt

	// Add role-specific data
	if user.Student != nil {
		u.Student = &StudentData{
			Instance: user.Student.Instance,
			Major:    user.Student.Major,
		}
	}

	if user.Mentor != nil {
		u.Mentor = &MentorData{
			Specialization: user.Mentor.Specialization,
			Experience:     user.Mentor.Experience,
			Rating:         user.Mentor.Rating,
			RatingCount:    user.Mentor.RatingCount,
			Price:          user.Mentor.Price,
		}
	}

	return u
}

func (u *UserResponse) PopulateMinimalFromEntity(user *entity.User) *UserResponse {
	u.ID = user.ID
	u.Name = user.Name
	u.Role = user.Role
	u.AvatarURL = user.AvatarURL

	if user.Student != nil {
		u.Student = &StudentData{
			Instance: user.Student.Instance,
			Major:    user.Student.Major,
		}
	}

	if user.Mentor != nil {
		u.Mentor = &MentorData{
			Specialization: user.Mentor.Specialization,
			Experience:     user.Mentor.Experience,
			Rating:         user.Mentor.Rating,
			RatingCount:    user.Mentor.RatingCount,
			Price:          user.Mentor.Price,
		}
	}

	return u
}

type UserUpdate struct {
	ID           uuid.UUID `db:"id"`
	Name         *string   `db:"name"`
	Email        *string   `db:"email"`
	PasswordHash *string   `db:"password_hash"`
	AvatarURL    *string   `db:"avatar_url"`

	Student *StudentUpdate `db:"student"`
	Mentor  *MentorUpdate  `db:"mentor"`
}

type StudentUpdate struct {
	Instance *string `db:"instance"`
	Major    *string `db:"major"`
}

type MentorUpdate struct {
	Specialization *string `db:"specialization"`
	Experience     *string `db:"experience"`
	Price          *int    `db:"price"`
}

type CreateStudentRequest struct {
	Instance string `form:"instance" json:"instance" validate:"required,min=1,max=50"`
	Major    string `form:"major" json:"major" validate:"required,min=1,max=50"`
}

type UpdateStudentRequest struct {
	Instance *string `form:"instance" json:"instance" validate:"omitempty,min=1,max=50"`
	Major    *string `form:"major" json:"major" validate:"omitempty,min=1,max=50"`
}

type CreateMentorRequest struct {
	Specialization string `form:"specialization" json:"specialization" validate:"required,min=1,max=255"`
	Experience     string `form:"experience" json:"experience" validate:"required,min=1,max=1000"`
	Price          int    `form:"price" json:"price" validate:"required,min=0"`
}

type UpdateMentorRequest struct {
	Specialization *string `form:"specialization" json:"specialization" validate:"omitempty,min=1,max=255"`
	Experience     *string `form:"experience" json:"experience" validate:"omitempty,min=1,max=1000"`
	Price          *int    `form:"price" json:"price" validate:"omitempty,min=0"`
}

type CreateUserRequest struct {
	Name     string                `json:"name" validate:"required,min=3,max=60"`
	Email    string                `json:"email" validate:"required,email,max=320"`
	Password string                `json:"-" validate:"required,min=8,max=72"`
	Role     enum.UserRole         `json:"role" validate:"required,oneof=admin mentor student"`
	Student  *CreateStudentRequest `json:"student,omitempty"`
	Mentor   *CreateMentorRequest  `json:"mentor,omitempty"`
}

type UpdateUserRequest struct {
	Name    *string               `form:"name" validate:"omitempty,min=3,max=60"`
	Avatar  *multipart.FileHeader `form:"avatar"`
	Student *UpdateStudentRequest `form:"student" json:"student,omitempty"`
	Mentor  *UpdateMentorRequest  `form:"mentor" json:"mentor,omitempty"`
}
