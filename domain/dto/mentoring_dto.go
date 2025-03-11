package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ChatResponse struct {
	ID        uuid.UUID `json:"id"`
	MentorID  uuid.UUID `json:"mentor_id"`
	StudentID uuid.UUID `json:"student_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (r *ChatResponse) PopulateFromEntity(chat *entity.MentoringChat) {
	r.ID = chat.ID
	r.MentorID = chat.MentorID
	r.StudentID = chat.StudentID
	r.ExpiresAt = chat.ExpiresAt
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	SenderID  uuid.UUID `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *MessageResponse) PopulateFromEntity(message *entity.MentoringMessage) {
	r.ID = message.ID
	r.SenderID = message.SenderID
	r.Message = message.Message
	r.CreatedAt = message.CreatedAt
}
