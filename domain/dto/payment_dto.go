package dto

import (
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type CreatePaymentRequest struct {
	UserID  uuid.UUID
	Amount  int
	Title   string
	Detail  *string
	Payload entity.PaymentPayload
}
