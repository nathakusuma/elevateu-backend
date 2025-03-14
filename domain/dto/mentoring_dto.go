package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type ChatResponse struct {
	ID                 uuid.UUID  `json:"id"`
	MentorID           uuid.UUID  `json:"mentor_id"`
	StudentID          uuid.UUID  `json:"student_id"`
	ExpiresAt          time.Time  `json:"expires_at"`
	IsTrial            bool       `json:"is_trial"`
	MentorAvatar       string     `json:"mentor_avatar,omitempty"`
	StudentAvatar      string     `json:"student_avatar,omitempty"`
	LastMessageContent string     `json:"last_message_content,omitempty"`
	LastMessageTime    *time.Time `json:"last_message_time,omitempty"`
}

func (r *ChatResponse) PopulateFromEntity(chat *entity.MentoringChat, urlSigner func(string) (string, error)) error {
	mentorAvatar, err := urlSigner("users/avatar/" + chat.MentorID.String())
	if err != nil {
		return err
	}

	studentAvatar, err := urlSigner("users/avatar/" + chat.StudentID.String())
	if err != nil {
		return err
	}

	r.ID = chat.ID
	r.MentorID = chat.MentorID
	r.StudentID = chat.StudentID
	r.ExpiresAt = chat.ExpiresAt
	r.IsTrial = chat.IsTrial
	r.MentorAvatar = mentorAvatar
	r.StudentAvatar = studentAvatar

	if chat.LastMessage != nil {
		r.LastMessageContent = chat.LastMessage.Message
		r.LastMessageTime = &chat.LastMessage.CreatedAt
	}

	return nil
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
