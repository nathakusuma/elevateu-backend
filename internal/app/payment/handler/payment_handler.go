package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
)

type paymentHandler struct {
	svc         contract.IPaymentService
	midtransSvc contract.IPaymentGateway
}

func InitPaymentHandler(
	router fiber.Router,
	svc contract.IPaymentService,
	midtransSvc contract.IPaymentGateway,
) {
	handler := paymentHandler{
		svc:         svc,
		midtransSvc: midtransSvc,
	}

	paymentGroup := router.Group("/payments")
	paymentGroup.Post("/midtrans/notification", handler.midtransNotification)
}

func (h *paymentHandler) midtransNotification(ctx *fiber.Ctx) error {
	var notificationPayload map[string]interface{}
	if err := ctx.BodyParser(&notificationPayload); err != nil {
		return errorpkg.ErrFailParseRequest
	}

	if err := h.midtransSvc.ProcessNotification(ctx.Context(),
		notificationPayload, h.svc.UpdatePaymentStatus); err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
