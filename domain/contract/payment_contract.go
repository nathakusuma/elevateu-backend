package contract

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type IPaymentRepository interface {
	CreatePayment(ctx context.Context, tx database.ITransaction, payment *entity.Payment) error
	CreateMentorTransactionHistory(ctx context.Context, txWrapper database.ITransaction,
		mentorTransactionHistory *entity.MentorTransactionHistory) error
	GetPaymentByID(ctx context.Context, tx database.ITransaction,
		id uuid.UUID) (*entity.Payment, error)
	UpdatePayment(ctx context.Context, tx database.ITransaction, payment *entity.Payment) error

	AddBoostSubscription(ctx context.Context, txWrapper database.ITransaction,
		studentID uuid.UUID, duration time.Duration) error
	AddChallengeSubscription(ctx context.Context, txWrapper database.ITransaction,
		studentID uuid.UUID, duration time.Duration) error
	AddMentorBalance(ctx context.Context, txWrapper database.ITransaction,
		mentorID uuid.UUID, amount int) error
}

type IPaymentService interface {
	UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status enum.PaymentStatus, method string) error
	ProcessNotification(ctx context.Context, notificationPayload map[string]any) error

	PaySkillBoost(ctx context.Context, studentID uuid.UUID) (string, error)
	PaySkillChallenge(ctx context.Context, studentID uuid.UUID) (string, error)
	PaySkillGuidance(ctx context.Context, studentID, mentorID uuid.UUID) (string, error)
}
