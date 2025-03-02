package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/nathakusuma/elevateu-backend/internal/infra/cache"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type paymentRepository struct {
	db    *sqlx.DB
	cache cache.ICache
}

func NewPaymentRepository(db *sqlx.DB, cache cache.ICache) contract.IPaymentRepository {
	return &paymentRepository{db: db, cache: cache}
}

func (r *paymentRepository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *paymentRepository) CreatePayment(ctx context.Context, tx sqlx.ExtContext, payment *entity.Payment,
	payload []*entity.PaymentPayload) error {
	if tx == nil {
		tx = r.db
	}

	if err := r.cache.Set(ctx, "payment:"+payment.ID.String(), payload, 1*time.Hour); err != nil {
		return fmt.Errorf("failed to cache payment payload: %w", err)
	}

	_, err := sqlx.NamedExecContext(ctx, tx, `
		INSERT INTO payments (
			id,
			user_id,
		    token,
			amount,
			title,
			detail,
			method,
			status,
		    expired_at
		) VALUES (
			:id,
			:user_id,
			:token,
			:amount,
			:title,
			:detail,
			:method,
			:status
			:expired_at
		)
	`, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *paymentRepository) GetPaymentByID(ctx context.Context, tx sqlx.QueryerContext,
	id uuid.UUID) (*entity.Payment, []*entity.PaymentPayload, error) {
	if tx == nil {
		tx = r.db
	}

	var payload []*entity.PaymentPayload
	if err := r.cache.Get(ctx, "payment:"+id.String(), &payload); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil, fmt.Errorf("payment payload not found: %w", err)
		}

		return nil, nil, fmt.Errorf("failed to get payment payload: %w", err)
	}

	var payment entity.Payment
	if err := sqlx.GetContext(ctx, tx, &payment, `
		SELECT
			id,
			user_id,
			token,
			amount,
			title,
			detail,
			method,
			status,
			expired_at,
			created_at,
			updated_at
		FROM payments
		WHERE id = $1
	`, id); err != nil {
		return nil, nil, fmt.Errorf("failed to get payment by ID: %w", err)
	}

	return &payment, payload, nil
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, tx sqlx.ExtContext, payment *entity.Payment) error {
	if tx == nil {
		tx = r.db
	}

	_, err := sqlx.NamedExecContext(ctx, tx, `
		UPDATE payments
		SET
			method = :method,
			status = :status
		WHERE id = :id
	`, payment)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}
