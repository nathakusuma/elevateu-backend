package contract

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type IPaymentGateway interface {
	CreateTransaction(id string, amount int) (string, error)
	ProcessNotification(ctx context.Context, notificationPayload map[string]interface{},
		statusUpdateCallback func(context.Context, uuid.UUID, enum.PaymentStatus) error) error
}

type IPaymentRepository interface {
	BeginTx() (*sqlx.Tx, error)
	CreatePayment(ctx context.Context, tx sqlx.ExtContext, payment *entity.Payment,
		payload []*entity.PaymentPayload) error
	GetPaymentByID(ctx context.Context, tx sqlx.QueryerContext,
		id uuid.UUID) (*entity.Payment, []*entity.PaymentPayload, error)
	UpdatePayment(ctx context.Context, tx sqlx.ExtContext, payment *entity.Payment) error
}

type IPaymentService interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error)
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status enum.PaymentStatus) error
}
