package dto

import (
	"github.com/google/uuid"
	"time"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type PaymentResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Token     string             `json:"token"`
	Amount    int                `json:"amount"`
	Title     string             `json:"title"`
	Detail    *string            `json:"detail"`
	Method    string             `json:"method"`
	Status    enum.PaymentStatus `json:"status"`
	ExpiredAt time.Time          `json:"expired_at"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func (p *PaymentResponse) PopulateFromEntity(payment *entity.Payment) {
	p.ID = payment.ID
	p.UserID = payment.UserID
	p.Token = payment.Token
	p.Amount = payment.Amount
	p.Title = payment.Title
	p.Detail = payment.Detail
	p.Method = payment.Method
	p.ExpiredAt = payment.ExpiredAt
	p.CreatedAt = payment.CreatedAt
	p.UpdatedAt = payment.UpdatedAt

	if payment.ExpiredAt.Before(time.Now()) {
		p.Status = enum.PaymentStatusFailure
	} else {
		p.Status = payment.Status
	}
}

type CreatePaymentRequest struct {
	UserID  uuid.UUID
	Amount  int
	Title   string
	Detail  *string
	Payload entity.PaymentPayload
}
