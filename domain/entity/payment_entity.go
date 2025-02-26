package entity

import (
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type Payment struct {
	ID        uuid.UUID          `db:"id"`
	UserID    uuid.UUID          `db:"user_id"`
	Token     string             `db:"token"`
	Amount    int                `db:"amount"`
	Title     string             `db:"title"`
	Detail    *string            `db:"detail"`
	Method    string             `db:"method"`
	Status    enum.PaymentStatus `db:"status"`
	ExpiredAt time.Time          `db:"expired_at"`
	CreatedAt time.Time          `db:"created_at"`
	UpdatedAt time.Time          `db:"updated_at"`
}

type PaymentPayload struct {
	Type enum.PaymentType
	Data any
}
