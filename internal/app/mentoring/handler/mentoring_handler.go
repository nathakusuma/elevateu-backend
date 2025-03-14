package handler

import (
	"context"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/jwt"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type mentoringHandler struct {
	svc contract.IMentoringService
	val validator.IValidator
	jwt jwt.IJwt
}

func InitMentoringHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	mentoringSvc contract.IMentoringService,
	jwt jwt.IJwt,
	validator validator.IValidator,
) {
	handler := mentoringHandler{
		val: validator,
		svc: mentoringSvc,
		jwt: jwt,
	}

	mentoringsGroup := router.Group("/mentorings")

	mentoringsGroup.Post("/trial",
		midw.RequireAuthenticated,
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.createTrialChat)

	mentoringsGroup.Get("/chats/my",
		midw.RequireAuthenticated,
		handler.getMyChats)

	mentoringsGroup.Post("/chats/:chatId/messages",
		midw.RequireAuthenticated,
		handler.sendMessage)
	mentoringsGroup.Get("/chats/:chatId/messages",
		midw.RequireAuthenticated,
		handler.getMessages)

	mentoringsGroup.Get("/chats/:chatId/ws",
		middleware.WebsocketUpgrade,
		websocket.New(handler.handleWebSocket),
	)
}

func (h *mentoringHandler) createTrialChat(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	var req struct {
		MentorID uuid.UUID `json:"mentor_id" validate:"required"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	chatResp, err := h.svc.CreateChat(ctx.Context(), req.MentorID, userID, true)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(chatResp)
}

func (h *mentoringHandler) getMyChats(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	chats, err := h.svc.GetChatsByUserID(ctx.Context(), userID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"chats": chats,
	})
}

func (h *mentoringHandler) sendMessage(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	chatID, err := uuid.Parse(ctx.Params("chatId"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid chat ID")
	}

	var req struct {
		Message string `json:"message" validate:"required,min=1,max=2000"`
	}
	if err = ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(req); err != nil {
		return err
	}

	if err = h.svc.SendMessage(ctx.Context(), userID, chatID, req.Message); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *mentoringHandler) getMessages(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		traceID := log.ErrorWithTraceID(ctx.Context(), nil, "Failed to get user ID from context")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	chatID, err := uuid.Parse(ctx.Params("chatId"))
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("Invalid chat ID")
	}

	var pageReq dto.PaginationRequest
	if err = ctx.QueryParser(&pageReq); err != nil {
		return errorpkg.ErrFailParseRequest()
	}

	if err = h.val.ValidateStruct(pageReq); err != nil {
		return err
	}

	messages, pageResp, err := h.svc.GetMessages(ctx.Context(), userID, chatID, pageReq)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"messages":   messages,
		"pagination": pageResp,
	})
}

func (h *mentoringHandler) handleWebSocket(conn *websocket.Conn) {
	token := conn.Headers("Authorization")
	if token == "" {
		conn.WriteJSON(errorpkg.ErrNoBearerToken)
		conn.Close()
		return
	}

	validateResp, err := h.jwt.Validate(token)
	if err != nil {
		conn.WriteJSON(errorpkg.ErrInvalidBearerToken)
		conn.Close()
		return
	}

	userID := validateResp.UserID

	chatIDStr := conn.Params("chatId")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		conn.WriteJSON(errorpkg.ErrValidation().WithDetail("Invalid chat ID"))
		conn.Close()
		return
	}

	// Set ping/pong handlers
	pingInterval := 30 * time.Second
	pongWait := 60 * time.Second

	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go func() {
		pingTicker := time.NewTicker(pingInterval)
		defer pingTicker.Stop()

		for range pingTicker.C {
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				if websocket.IsCloseError(err,
					websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure) {
					return
				}

				log.Error(context.Background(), map[string]interface{}{
					"error":   err,
					"user.id": userID,
					"chat.id": chatID,
				}, "Failed to send ping")
				return
			}
		}
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))

	h.svc.RegisterClient(userID, chatID, conn)

	var (
		msg []byte
	)

	defer func() {
		h.svc.UnregisterClient(userID, chatID)
		conn.Close()
	}()

	for {
		if _, msg, err = conn.ReadMessage(); err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure) {
				break
			}

			log.Error(context.Background(), map[string]interface{}{
				"error":   err,
				"user.id": userID,
				"chat.id": chatID,
			}, "WebSocket read error")
			break
		}

		if err = h.svc.SendMessage(context.Background(), userID, chatID, string(msg)); err != nil {
			conn.WriteJSON(err)
			break
		}
	}
}
