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
	AvatarURL    *string       `db:"avatar_url"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`

	Student *Student `db:"student"`
	Mentor  *Mentor  `db:"mentor"`
}

type Student struct {
	Instance string `db:"instance"`
	Major    string `db:"major"`
}

type Mentor struct {
	Specialization string  `db:"specialization"`
	Experience     string  `db:"experience"`
	Rating         float64 `db:"rating"`
	RatingCount    int     `db:"rating_count"`
	RatingTotal    float64 `db:"rating_total"`
	Price          int     `db:"price"`
	Balance        int     `db:"balance"`
}
