package entity

import (
	"time"

	"github.com/google/uuid"
)

type MentoringChat struct {
	ID        uuid.UUID `db:"id"`
	MentorID  uuid.UUID `db:"mentor_id"`
	StudentID uuid.UUID `db:"student_id"`
	ExpiresAt time.Time `db:"expires_at"`
}

type MentoringMessage struct {
	ID        uuid.UUID `db:"id"`
	ChatID    uuid.UUID `db:"chat_id"`
	SenderID  uuid.UUID `db:"sender_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
