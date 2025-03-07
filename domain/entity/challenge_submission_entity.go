package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChallengeSubmission struct {
	ID          uuid.UUID `db:"id"`
	ChallengeID uuid.UUID `db:"challenge_id"`
	StudentID   uuid.UUID `db:"student_id"`
	URL         string    `db:"url"`
	CreatedAt   time.Time `db:"created_at"`

	Student  *User                        `db:"student"`
	Feedback *ChallengeSubmissionFeedback `db:"feedback"`
}

type ChallengeSubmissionFeedback struct {
	SubmissionID uuid.UUID `db:"submission_id"`
	MentorID     uuid.UUID `db:"mentor_id"`
	Score        int       `db:"score"`
	Feedback     string    `db:"feedback"`
	CreatedAt    time.Time `db:"created_at"`

	Mentor *User `db:"mentor"`
}
