package handler

import (
	"context"

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
		return errorpkg.ErrInvalidBearerToken()
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

func (h *mentoringHandler) sendMessage(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken()
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
		return errorpkg.ErrInvalidBearerToken()
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
	token := conn.Query("token")
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
			conn.WriteJSON(errorpkg.ErrFailReadMessage)
			break
		}

		if err = h.svc.SendMessage(context.Background(), userID, chatID, string(msg)); err != nil {
			conn.WriteJSON(err)
			break
		}
	}
}
