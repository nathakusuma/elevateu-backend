package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) contract.IPaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreatePayment(ctx context.Context, txWrapper database.ITransaction,
	payment *entity.Payment) error {
	tx := txWrapper.GetTx()

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
			:status,
			:expired_at
		)
	`, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *paymentRepository) CreateMentorTransactionHistory(ctx context.Context, txWrapper database.ITransaction,
	mentorTransactionHistory *entity.MentorTransactionHistory) error {
	tx := txWrapper.GetTx()

	_, err := sqlx.NamedExecContext(ctx, tx, `
		INSERT INTO mentor_transaction_histories (
			id,
			mentor_id,
			title,
			detail,
			amount
		) VALUES (
			:id,
			:mentor_id,
			:title,
			:detail,
			:amount
		)
	`, mentorTransactionHistory)
	if err != nil {
		return fmt.Errorf("failed to create mentor transaction history: %w", err)
	}

	return nil
}

func (r *paymentRepository) GetPaymentByID(ctx context.Context, txWrapper database.ITransaction,
	id uuid.UUID) (*entity.Payment, error) {
	tx := txWrapper.GetTx()

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
		return nil, fmt.Errorf("failed to get payment by ID: %w", err)
	}

	return &payment, nil
}

func (r *paymentRepository) GetPaymentsByStudent(ctx context.Context, studentID uuid.UUID) ([]*entity.Payment, error) {
	var payments []*entity.Payment
	if err := sqlx.SelectContext(ctx, r.db, &payments, `
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
		WHERE user_id = $1
		ORDER BY id DESC
	`, studentID); err != nil {
		return nil, fmt.Errorf("failed to get payments by student: %w", err)
	}

	return payments, nil
}

func (r *paymentRepository) GetTransactionHistoriesByMentor(ctx context.Context, mentorID uuid.UUID) (
	[]*entity.MentorTransactionHistory, error) {
	var mentorTransactionHistories []*entity.MentorTransactionHistory
	if err := sqlx.SelectContext(ctx, r.db, &mentorTransactionHistories, `
		SELECT
			id,
			mentor_id,
			title,
			detail,
			amount,
			created_at
		FROM mentor_transaction_histories
		WHERE mentor_id = $1
		ORDER BY id DESC
	`, mentorID); err != nil {
		return nil, fmt.Errorf("failed to get transaction histories by mentor: %w", err)
	}

	return mentorTransactionHistories, nil
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, txWrapper database.ITransaction,
	payment *entity.Payment) error {
	tx := txWrapper.GetTx()

	_, err := sqlx.NamedExecContext(ctx, tx, `
		UPDATE payments
		SET
			method = :method,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
	`, payment)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

func (r *paymentRepository) AddBoostSubscription(ctx context.Context, txWrapper database.ITransaction,
	studentID uuid.UUID, subscribedUntil time.Time) error {
	tx := txWrapper.GetTx()

	_, err := tx.ExecContext(ctx, `
		UPDATE students SET subscribed_boost_until = $1
		WHERE user_id = $2
	`, subscribedUntil, studentID)

	if err != nil {
		return fmt.Errorf("failed to add skill boost subscription: %w", err)
	}

	return nil
}

func (r *paymentRepository) AddChallengeSubscription(ctx context.Context, txWrapper database.ITransaction,
	studentID uuid.UUID, subscribedUntil time.Time) error {
	tx := txWrapper.GetTx()

	_, err := tx.ExecContext(ctx, `
		UPDATE students SET subscribed_challenge_until = $1
		WHERE user_id = $2
	`, subscribedUntil, studentID)

	if err != nil {
		return fmt.Errorf("failed to add challenge subscription: %w", err)
	}

	return nil
}

func (r *paymentRepository) AddMentorBalance(ctx context.Context, txWrapper database.ITransaction,
	mentorID uuid.UUID, amount int) error {
	tx := txWrapper.GetTx()

	_, err := tx.ExecContext(ctx, `
		UPDATE mentors SET balance = balance + $1
		WHERE user_id = $2
	`, amount, mentorID)

	if err != nil {
		return fmt.Errorf("failed to add mentor balance: %w", err)
	}

	return nil
}
