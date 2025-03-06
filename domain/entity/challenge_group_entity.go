package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChallengeGroup struct {
	ID             uuid.UUID `db:"id"`
	CategoryID     uuid.UUID `db:"category_id"`
	Title          string    `db:"title"`
	Description    string    `db:"description"`
	ChallengeCount int       `db:"challenge_count"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`

	Category *Category `db:"category"`
}
