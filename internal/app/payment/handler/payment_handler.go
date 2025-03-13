package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/middleware"
	"github.com/nathakusuma/elevateu-backend/pkg/validator"
)

type paymentHandler struct {
	svc  contract.IPaymentService
	midw *middleware.Middleware
	val  validator.IValidator
}

func InitPaymentHandler(
	router fiber.Router,
	midw *middleware.Middleware,
	svc contract.IPaymentService,
	val validator.IValidator,
) {
	handler := paymentHandler{
		svc: svc,
		val: val,
	}

	paymentGroup := router.Group("/payments")
	paymentGroup.Post("/midtrans/notifications", handler.midtransNotification)

	paymentGroup.Post("/skill-boost",
		midw.RequireAuthenticated,
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.paySkillBoost)
	paymentGroup.Post("/skill-challenge",
		midw.RequireAuthenticated,
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.paySkillChallenge)
	paymentGroup.Post("/skill-guidance",
		midw.RequireAuthenticated,
		midw.RequireOneOfRoles(enum.UserRoleStudent),
		handler.paySkillGuidance)
}

func (h *paymentHandler) midtransNotification(ctx *fiber.Ctx) error {
	var notificationPayload map[string]interface{}
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.svc.ProcessNotification(ctx.Context(), notificationPayload); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *paymentHandler) paySkillBoost(ctx *fiber.Ctx) error {
	studentID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	paymentToken, err := h.svc.PaySkillBoost(ctx.Context(), studentID)
	if err != nil {
		return err
	}

	return ctx.JSON(map[string]any{
		"payment_token": paymentToken,
	})
}

func (h *paymentHandler) paySkillChallenge(ctx *fiber.Ctx) error {
	studentID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	paymentToken, err := h.svc.PaySkillChallenge(ctx.Context(), studentID)
	if err != nil {
		return err
	}

	return ctx.JSON(map[string]any{
		"payment_token": paymentToken,
	})
}

func (h *paymentHandler) paySkillGuidance(ctx *fiber.Ctx) error {
	studentID, ok := ctx.Locals(ctxkey.UserID).(uuid.UUID)
	if !ok {
		return errorpkg.ErrInvalidBearerToken
	}

	var req struct {
		MentorID uuid.UUID `json:"mentor_id" validate:"required"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return err
	}

	paymentToken, err := h.svc.PaySkillGuidance(ctx.Context(), studentID, req.MentorID)
	if err != nil {
		return err
	}

	return ctx.JSON(map[string]any{
		"payment_token": paymentToken,
	})
}
