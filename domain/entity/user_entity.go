package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type User struct {
	ID           uuid.UUID     `db:"id"`
	Name         string        `db:"name"`
	Email        string        `db:"email"`
	PasswordHash string        `db:"password_hash"`
	Role         enum.UserRole `db:"role"`
	Bio          *string       `db:"bio"`
	AvatarURL    *string       `db:"avatar_url"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
	DeletedAt    *time.Time    `db:"deleted_at"`
}
