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
	AvatarURL *string       `json:"avatar_url,omitempty"`
	CreatedAt *time.Time    `json:"created_at,omitempty"`
	UpdatedAt *time.Time    `json:"updated_at,omitempty"`
	Student   *StudentData  `json:"student,omitempty"`
	Mentor    *MentorData   `json:"mentor,omitempty"`
}

type StudentData struct {
	Instance                 string            `json:"instance,omitempty"`
	Major                    string            `json:"major,omitempty"`
	Point                    *int              `json:"point,omitempty"`
	Badge                    enum.StudentBadge `json:"badge,omitempty"`
	SubscribedBoostUntil     *time.Time        `json:"subscribed_boost_until,omitempty"`
	SubscribedChallengeUntil *time.Time        `json:"subscribed_challenge_until,omitempty"`
}

type MentorData struct {
	Address        string  `json:"address,omitempty"`
	Specialization string  `json:"specialization,omitempty"`
	CurrentJob     string  `json:"current_job,omitempty"`
	Company        string  `json:"company,omitempty"`
	Bio            *string `json:"bio,omitempty"`
	Gender         string  `json:"gender,omitempty"`
	Rating         float64 `json:"rating,omitempty"`
	RatingCount    int     `json:"rating_count,omitempty"`
	Price          int     `json:"price,omitempty"`
	Balance        int     `json:"balance,omitempty"`
}

func (u *UserResponse) PopulateFromEntity(user *entity.User,
	urlSigner func(string) (string, error)) error {
	u.ID = user.ID
	u.Name = user.Name
	u.Email = user.Email
	u.Role = user.Role
	u.CreatedAt = &user.CreatedAt
	u.UpdatedAt = &user.UpdatedAt

	var avatarURL string
	var err error
	if user.HasAvatar {
		avatarURL, err = urlSigner("users/avatar/" + user.ID.String())
	} else {
		avatarURL, err = urlSigner("users/avatar/default")
	}
	if err != nil {
		return err
	}
	u.AvatarURL = &avatarURL

	// Add role-specific data
	if user.Student != nil {
		u.Student = &StudentData{
			Instance: user.Student.Instance,
			Major:    user.Student.Major,
			Point:    &user.Student.Point,
		}
		u.Student.Badge = enum.GetBadge(user.Student.Point)
		if !user.Student.SubscribedBoostUntil.IsZero() {
			u.Student.SubscribedBoostUntil = &user.Student.SubscribedBoostUntil
		}
		if !user.Student.SubscribedChallengeUntil.IsZero() {
			u.Student.SubscribedChallengeUntil = &user.Student.SubscribedChallengeUntil
		}
	}

	if user.Mentor != nil {
		u.Mentor = &MentorData{
			Address:        user.Mentor.Address,
			Specialization: user.Mentor.Specialization,
			CurrentJob:     user.Mentor.CurrentJob,
			Company:        user.Mentor.Company,
			Bio:            user.Mentor.Bio,
			Gender:         user.Mentor.Gender,
			Rating:         user.Mentor.Rating,
			RatingCount:    user.Mentor.RatingCount,
			Price:          user.Mentor.Price,
			Balance:        user.Mentor.Balance,
		}
	}

	return nil
}

func (u *UserResponse) PopulateMinimalFromEntity(user *entity.User,
	urlSigner func(string) (string, error)) error {
	u.ID = user.ID
	u.Name = user.Name
	u.Role = user.Role

	if user.HasAvatar {
		avatarURL, err := urlSigner("users/avatar/" + user.ID.String())
		if err != nil {
			return err
		}
		u.AvatarURL = &avatarURL
	}

	// Add role-specific data
	if user.Student != nil {
		u.Student = &StudentData{
			Instance: user.Student.Instance,
			Major:    user.Student.Major,
			Point:    &user.Student.Point,
			Badge:    enum.GetBadge(user.Student.Point),
		}
	}

	if user.Mentor != nil {
		u.Mentor = &MentorData{
			Specialization: user.Mentor.Specialization,
			CurrentJob:     user.Mentor.CurrentJob,
			Company:        user.Mentor.Company,
			Bio:            user.Mentor.Bio,
			Gender:         user.Mentor.Gender,
			Rating:         user.Mentor.Rating,
			RatingCount:    user.Mentor.RatingCount,
			Price:          user.Mentor.Price,
		}
	}

	return nil
}

type UserUpdate struct {
	ID           uuid.UUID `db:"id"`
	Name         *string   `db:"name"`
	Email        *string   `db:"email"`
	PasswordHash *string   `db:"password_hash"`
	HasAvatar    *bool     `db:"has_avatar"`

	Student *StudentUpdate `db:"-"`
	Mentor  *MentorUpdate  `db:"-"`
}

type StudentUpdate struct {
	Instance *string `db:"instance"`
	Major    *string `db:"major"`
}

type MentorUpdate struct {
	Address        *string `db:"address"`
	Specialization *string `db:"specialization"`
	CurrentJob     *string `db:"current_job"`
	Company        *string `db:"company"`
	Bio            *string `db:"bio"`
	Gender         *string `db:"gender"`
	Price          *int    `db:"price"`
}

type CreateUserRequest struct {
	Name     string                `json:"name"`
	Email    string                `json:"email"`
	Password string                `json:"-"`
	Role     enum.UserRole         `json:"role"`
	Student  *CreateStudentRequest `json:"student,omitempty"`
	Mentor   *CreateMentorRequest  `json:"mentor,omitempty"`
}

type CreateStudentRequest struct {
	Instance string `json:"instance" validate:"required,min=1,max=50"`
	Major    string `json:"major" validate:"required,min=1,max=50"`
}

type CreateMentorRequest struct {
	Address        string `json:"address" validate:"required,min=1,max=255"`
	Specialization string `json:"specialization" validate:"required,min=1,max=50"`
	CurrentJob     string `json:"current_job" validate:"required,min=1,max=50"`
	Company        string `json:"company" validate:"required,min=1,max=50"`
	Gender         string `json:"gender" validate:"required,oneof=male female"`
}

type UpdateUserRequest struct {
	Name    *string               `json:"name" validate:"omitempty,min=3,max=60"`
	Student *UpdateStudentRequest `json:"student,omitempty"`
	Mentor  *UpdateMentorRequest  `json:"mentor,omitempty"`
}

type UpdateStudentRequest struct {
	Instance *string `json:"instance" validate:"omitempty,min=1,max=50"`
	Major    *string `json:"major" validate:"omitempty,min=1,max=50"`
}

type UpdateMentorRequest struct {
	Address        *string `json:"address" validate:"omitempty,min=1,max=255"`
	Specialization *string `json:"specialization" validate:"omitempty,min=1,max=255"`
	CurrentJob     *string `json:"current_job" validate:"omitempty,min=1,max=255"`
	Company        *string `json:"company" validate:"omitempty,min=1,max=255"`
	Bio            *string `json:"bio" validate:"omitempty,min=1,max=255"`
	Gender         *string `json:"gender" validate:"omitempty,oneof=male female"`
}
