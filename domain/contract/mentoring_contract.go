package contract

import (
	"context"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type IMentoringRepository interface {
	CreateChat(ctx context.Context, chat *entity.MentoringChat) error
	CreateTrialChat(ctx context.Context, chat *entity.MentoringChat) error
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*entity.MentoringChat, error)
	SendMessage(ctx context.Context, message *entity.MentoringMessage) error
	GetMessages(ctx context.Context, chatID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*entity.MentoringMessage, dto.PaginationResponse, error)
}

type IMentoringService interface {
	CreateChat(ctx context.Context, mentorID, studentID uuid.UUID, isTrial bool) (*dto.ChatResponse, error)
	SendMessage(ctx context.Context, userID, chatID uuid.UUID,
		message string) error
	GetMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID,
		pageReq dto.PaginationRequest) ([]*dto.MessageResponse, dto.PaginationResponse, error)

	RegisterClient(userID uuid.UUID, chatID uuid.UUID, conn *websocket.Conn)
	UnregisterClient(userID uuid.UUID, chatID uuid.UUID)
	BroadcastMessage(message *dto.MessageResponse, chatID uuid.UUID)
}
