package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type ChallengeResponse struct {
	ID              uuid.UUID                `json:"id"`
	Title           string                   `json:"title,omitempty"`
	Subtitle        string                   `json:"subtitle,omitempty"`
	Description     string                   `json:"description,omitempty"`
	Difficulty      enum.ChallengeDifficulty `json:"difficulty,omitempty"`
	IsFree          *bool                    `json:"is_free,omitempty"`
	SubmissionCount int64                    `json:"submission_count,omitempty"`
	CreatedAt       *time.Time               `json:"created_at,omitempty"`
	UpdatedAt       *time.Time               `json:"updated_at,omitempty"`

	SubmissionURL string `json:"submission_url,omitempty"`
}

func (r *ChallengeResponse) PopulateFromEntity(challenge *entity.Challenge) {
	r.ID = challenge.ID
	r.Title = challenge.Title
	r.Subtitle = challenge.Subtitle
	r.Difficulty = challenge.Difficulty
	r.IsFree = &challenge.IsFree
	r.SubmissionCount = challenge.SubmissionCount

	if !challenge.CreatedAt.IsZero() {
		r.CreatedAt = &challenge.CreatedAt
	}
	if !challenge.UpdatedAt.IsZero() {
		r.UpdatedAt = &challenge.UpdatedAt
	}
}

func (r *ChallengeResponse) PopulateDetailFromEntity(challenge *entity.Challenge) {
	r.PopulateFromEntity(challenge)
	r.Description = challenge.Description
	if challenge.Submission != nil {
		r.SubmissionURL = challenge.Submission.URL
	}
}

type ChallengeUpdate struct {
	GroupID     *uuid.UUID                `db:"group_id"`
	Title       *string                   `db:"title"`
	Subtitle    *string                   `db:"subtitle"`
	Description *string                   `db:"description"`
	Difficulty  *enum.ChallengeDifficulty `db:"difficulty"`
	IsFree      *bool                     `db:"is_free"`
}

type CreateChallengeRequest struct {
	GroupID     uuid.UUID                `json:"group_id" validate:"required"`
	Title       string                   `json:"title" validate:"required,min=3,max=50"`
	Subtitle    string                   `json:"subtitle" validate:"required,min=3,max=100"`
	Description string                   `json:"description" validate:"required,min=3,max=5000"`
	Difficulty  enum.ChallengeDifficulty `json:"difficulty" validate:"required,oneof=beginner intermediate advanced"`
	IsFree      bool                     `json:"is_free"`
}

type UpdateChallengeRequest struct {
	GroupID     *uuid.UUID                `json:"group_id" validate:"omitempty"`
	Title       *string                   `json:"title" validate:"omitempty,min=3,max=50"`
	Subtitle    *string                   `json:"subtitle" validate:"omitempty,min=3,max=100"`
	Description *string                   `json:"description" validate:"omitempty,min=3,max=5000"`
	Difficulty  *enum.ChallengeDifficulty `json:"difficulty" validate:"omitempty,oneof=beginner intermediate advanced"`
	IsFree      *bool                     `json:"is_free"`
}
