package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type mentoringService struct {
	repo         contract.IMentoringRepository
	userRepo     contract.IUserRepository
	uuid         uuidpkg.IUUID
	clients      map[string]map[uuid.UUID]*websocket.Conn
	clientsMutex sync.RWMutex
}

func NewMentoringService(
	mentoringRepo contract.IMentoringRepository,
	userRepo contract.IUserRepository,
	uuidGen uuidpkg.IUUID,
) contract.IMentoringService {
	return &mentoringService{
		repo:     mentoringRepo,
		userRepo: userRepo,
		uuid:     uuidGen,
		clients:  make(map[string]map[uuid.UUID]*websocket.Conn),
	}
}

func (s *mentoringService) CreateChat(ctx context.Context, mentorID,
	studentID uuid.UUID, isTrial bool) (*dto.ChatResponse, error) {
	mentor, err := s.userRepo.GetUserByField(ctx, "id", mentorID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "user not found") {
			return nil, errorpkg.ErrValidation().WithDetail("Mentor not found")
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"mentor.id": mentorID,
		}, "Failed to get mentor")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if mentor.Role != enum.UserRoleMentor {
		return nil, errorpkg.ErrValidation().WithDetail("User is not a mentor")
	}

	student, err := s.userRepo.GetUserByField(ctx, "id", studentID)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"mentor.id": mentorID,
		}, "Failed to get student")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if student.Role != enum.UserRoleStudent {
		return nil, errorpkg.ErrValidation().WithDetail("User is not a student")
	}

	currentChat, err := s.repo.GetChatByMentorAndStudent(ctx, mentorID, studentID)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "chat not found") {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":     err,
				"mentor.id": mentorID,
			}, "Failed to get chat")
			return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
		}
	}

	if isTrial && currentChat != nil {
		return nil, errorpkg.ErrTrialUsed()
	}

	chatID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":     err,
			"mentor.id": mentorID,
		}, "Failed to generate chat ID")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	expireDuration := 24 * time.Hour
	if isTrial {
		expireDuration = 15 * time.Minute
	}

	var expiresAt time.Time
	if currentChat != nil && currentChat.ExpiresAt.After(time.Now()) {
		expiresAt = currentChat.ExpiresAt.Add(expireDuration)
	} else {
		expiresAt = time.Now().Add(expireDuration)
	}

	chat := &entity.MentoringChat{
		ID:        chatID,
		MentorID:  mentorID,
		StudentID: studentID,
		ExpiresAt: expiresAt,
		IsTrial:   isTrial,
	}

	var repoErr error
	if isTrial {
		repoErr = s.repo.CreateTrialChat(ctx, chat)
	} else {
		repoErr = s.repo.CreateChat(ctx, chat)
	}

	if repoErr != nil {
		if strings.HasPrefix(repoErr.Error(), "trial chat already exists") {
			return nil, errorpkg.ErrTrialUsed()
		}
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": repoErr,
			"chat":  chat,
		}, "Failed to create chat")
		return nil, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	response := &dto.ChatResponse{}
	response.PopulateFromEntity(chat)

	log.Info(ctx, map[string]interface{}{
		"mentor.id": mentorID,
		"chat.id":   chatID,
	}, "Chat created")

	return response, nil
}

func (s *mentoringService) SendMessage(ctx context.Context, userID, chatID uuid.UUID,
	message string) error {
	chat, err := s.repo.GetChatByID(ctx, chatID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "chat not found") {
			return errorpkg.ErrValidation().WithDetail("Chat not found")
		}
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"chat.id": chatID,
		}, "Failed to verify chat access")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	isParticipant := chat.MentorID == userID || chat.StudentID == userID
	if !isParticipant {
		return errorpkg.ErrForbiddenUser().WithDetail("You don't have access to this chat")
	}

	if time.Now().After(chat.ExpiresAt) {
		return errorpkg.ErrChatExpired()
	}

	messageID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"chat.id": chatID,
		}, "Failed to generate message ID")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	messageEntity := &entity.MentoringMessage{
		ID:       messageID,
		ChatID:   chatID,
		SenderID: userID,
		Message:  message,
	}

	if err = s.repo.SendMessage(ctx, messageEntity); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error": err,
		}, "Failed to send message")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	response := &dto.MessageResponse{}
	response.PopulateFromEntity(messageEntity)

	s.BroadcastMessage(response, chatID)

	return nil
}

func (s *mentoringService) GetMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*dto.MessageResponse, dto.PaginationResponse, error) {

	chat, err := s.repo.GetChatByID(ctx, chatID)
	if err != nil {
		if strings.HasPrefix(err.Error(), "chat not found") {
			return nil, dto.PaginationResponse{},
				errorpkg.ErrValidation().WithDetail("Chat not found")
		}
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"chat.id": chatID,
		}, "Failed to verify chat access")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	isParticipant := chat.MentorID == userID || chat.StudentID == userID
	if !isParticipant {
		return nil, dto.PaginationResponse{},
			errorpkg.ErrForbiddenUser().WithDetail("You don't have access to this chat")
	}

	messages, pageResp, err := s.repo.GetMessages(ctx, chatID, pageReq)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"chat.id": chatID,
		}, "Failed to get messages")
		return nil, dto.PaginationResponse{}, errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	responses := make([]*dto.MessageResponse, len(messages))
	for i, message := range messages {
		responses[i] = &dto.MessageResponse{}
		responses[i].PopulateFromEntity(message)
	}

	return responses, pageResp, nil
}

func (s *mentoringService) RegisterClient(userID uuid.UUID, chatID uuid.UUID, conn *websocket.Conn) {
	chatIDStr := chatID.String()

	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if _, ok := s.clients[chatIDStr]; !ok {
		s.clients[chatIDStr] = make(map[uuid.UUID]*websocket.Conn)
	}

	s.clients[chatIDStr][userID] = conn

	log.Info(context.Background(), map[string]interface{}{
		"user.id": userID,
		"chat.id": chatID,
	}, "Client registered for chat")
}

func (s *mentoringService) UnregisterClient(userID uuid.UUID, chatID uuid.UUID) {
	chatIDStr := chatID.String()

	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()

	if chatClients, ok := s.clients[chatIDStr]; ok {
		delete(chatClients, userID)

		if len(chatClients) == 0 {
			delete(s.clients, chatIDStr)
		}

		log.Info(context.Background(), map[string]interface{}{
			"user.id": userID,
			"chat.id": chatID,
		}, "Client unregistered from chat")
	}
}

func (s *mentoringService) BroadcastMessage(message *dto.MessageResponse, chatID uuid.UUID) {
	chatIDStr := chatID.String()

	s.clientsMutex.RLock()
	defer s.clientsMutex.RUnlock()

	chatClients, ok := s.clients[chatIDStr]
	if !ok {
		return
	}

	for userID, conn := range chatClients {
		if userID == message.SenderID {
			continue
		}

		if err := conn.WriteJSON(message); err != nil {
			log.Error(context.Background(), map[string]interface{}{
				"error":   err,
				"user.id": userID,
				"chat.id": chatID,
			}, "Failed to send message to client")
		}
	}
}
