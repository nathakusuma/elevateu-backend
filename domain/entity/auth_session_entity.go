package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	Token     string    `db:"token"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
