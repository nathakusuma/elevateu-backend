package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type Challenge struct {
	ID              uuid.UUID                `db:"id"`
	GroupID         uuid.UUID                `db:"group_id"`
	Title           string                   `db:"title"`
	Subtitle        string                   `db:"subtitle"`
	Description     string                   `db:"description"`
	Difficulty      enum.ChallengeDifficulty `db:"difficulty"`
	IsFree          bool                     `db:"is_free"`
	SubmissionCount int64                    `db:"submission_count"`
	CreatedAt       time.Time                `db:"created_at"`
	UpdatedAt       time.Time                `db:"updated_at"`

	Submission *ChallengeSubmission `db:"submission"`
}
