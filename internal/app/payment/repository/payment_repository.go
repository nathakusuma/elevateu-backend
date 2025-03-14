package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
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

func (r *paymentRepository) GetPaymentsByStudent(ctx context.Context, studentID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.Payment, dto.PaginationResponse, error) {
	baseQuery := `
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
	`

	var sqlQuery string
	var args []interface{}
	args = append(args, studentID)

	if pageReq.Cursor != uuid.Nil {
		var operator string
		var orderDirection string

		if pageReq.Direction == "next" {
			operator = "<"
			orderDirection = "DESC"
		} else {
			operator = ">"
			orderDirection = "ASC"
		}

		sqlQuery = baseQuery + fmt.Sprintf(" AND id %s $2 ORDER BY id %s LIMIT $3", operator, orderDirection)
		args = append(args, pageReq.Cursor, pageReq.Limit+1)
	} else {
		sqlQuery = baseQuery + " ORDER BY id DESC LIMIT $2"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get payments: %w", err)
	}
	defer rows.Close()

	var payments []*entity.Payment
	for rows.Next() {
		var payment entity.Payment
		if err := rows.StructScan(&payment); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan payment row: %w", err)
		}
		payments = append(payments, &payment)
	}

	if err := rows.Err(); err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over payment rows: %w", err)
	}

	hasMore := false
	if len(payments) > pageReq.Limit {
		hasMore = true
		payments = payments[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		for i, j := 0, len(payments)-1; i < j; i, j = i+1, j-1 {
			payments[i], payments[j] = payments[j], payments[i]
		}
	}

	return payments, dto.PaginationResponse{HasMore: hasMore}, nil
}

func (r *paymentRepository) GetTransactionHistoriesByMentor(ctx context.Context, mentorID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.MentorTransactionHistory, dto.PaginationResponse, error) {
	baseQuery := `
		SELECT
			id,
			mentor_id,
			title,
			detail,
			amount,
			created_at
		FROM mentor_transaction_histories
		WHERE mentor_id = $1
	`

	var sqlQuery string
	var args []interface{}
	args = append(args, mentorID)

	if pageReq.Cursor != uuid.Nil {
		var operator string
		var orderDirection string

		if pageReq.Direction == "next" {
			operator = "<"
			orderDirection = "DESC"
		} else {
			operator = ">"
			orderDirection = "ASC"
		}

		sqlQuery = baseQuery + fmt.Sprintf(" AND id %s $2 ORDER BY id %s LIMIT $3", operator, orderDirection)
		args = append(args, pageReq.Cursor, pageReq.Limit+1)
	} else {
		sqlQuery = baseQuery + " ORDER BY id DESC LIMIT $2"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get transaction histories: %w", err)
	}
	defer rows.Close()

	var histories []*entity.MentorTransactionHistory
	for rows.Next() {
		var history entity.MentorTransactionHistory
		if err := rows.StructScan(&history); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan transaction history row: %w", err)
		}
		histories = append(histories, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over transaction history rows: %w", err)
	}

	hasMore := false
	if len(histories) > pageReq.Limit {
		hasMore = true
		histories = histories[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		for i, j := 0, len(histories)-1; i < j; i, j = i+1, j-1 {
			histories[i], histories[j] = histories[j], histories[i]
		}
	}

	return histories, dto.PaginationResponse{HasMore: hasMore}, nil
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
