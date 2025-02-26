package dto

import "github.com/nathakusuma/elevateu-backend/domain/enum"

type RegisterRequest struct {
	Email    string                `json:"email"    validate:"required,email,max=320"`
	OTP      string                `json:"otp"      validate:"required"`
	Name     string                `json:"name"     validate:"required,min=3,max=60,ascii"`
	Password string                `json:"password" validate:"required,min=8,max=72,ascii"`
	Role     enum.UserRole         `json:"role"     validate:"required,oneof=student mentor"`
	Student  *CreateStudentRequest `json:"student,omitempty" validate:"required_if=Role student"`
	Mentor   *CreateMentorRequest  `json:"mentor,omitempty" validate:"required_if=Role mentor"`
}

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,ascii"`
}

type LoginResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"        validate:"required,email"`
	OTP         string `json:"otp"          validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72,ascii"`
}
