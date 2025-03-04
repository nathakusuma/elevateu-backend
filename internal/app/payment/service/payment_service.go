package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type paymentService struct {
	repo       contract.IPaymentRepository
	gatewaySvc contract.IPaymentGateway
	uuid       uuidpkg.IUUID
}

func NewPaymentService(
	repo contract.IPaymentRepository,
	gatewaySvc contract.IPaymentGateway,
	uuid uuidpkg.IUUID,
) contract.IPaymentService {
	return &paymentService{
		repo:       repo,
		gatewaySvc: gatewaySvc,
		uuid:       uuid,
	}
}

func (s *paymentService) CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error) {
	paymentID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[PaymentService][CreatePayment] Failed to generate payment ID")
		return "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	// Create transaction in payment gateway
	token, err := s.gatewaySvc.CreateTransaction(paymentID.String(), req.Amount)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err,
			"request": req,
		}, "[PaymentService][CreatePayment] Failed to create transaction in payment gateway")
		return "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	payment := &entity.Payment{
		ID:        paymentID,
		UserID:    req.UserID,
		Token:     token,
		Amount:    req.Amount,
		Title:     req.Title,
		Detail:    req.Detail,
		Status:    enum.PaymentStatusPending,
		ExpiredAt: time.Now().Add(1 * time.Hour),
	}

	if err2 := s.repo.CreatePayment(ctx, nil, payment, req.Payload); err2 != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":   err2,
			"request": req,
		}, "[PaymentService][CreatePayment] Failed to create payment")
		return "", errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return token, nil
}

func (s *paymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status enum.PaymentStatus) error {
	tx, err := s.repo.BeginTx()
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":          err,
			"payment.id":     id,
			"payment.status": status,
		}, "[PaymentService][UpdatePaymentStatus] Failed to begin transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}
	defer func() {
		if err2 := tx.Rollback(); err2 != nil && !errors.Is(err2, sql.ErrTxDone) {
			log.Error(map[string]interface{}{
				"error":          err2,
				"payment.id":     id,
				"payment.status": status,
			}, "[PaymentService][UpdatePaymentStatus] Failed to rollback transaction")
		}
	}()

	payment, payload, err := s.repo.GetPaymentByID(ctx, tx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "payment payload not found") {
			return errorpkg.ErrNotFound
		}

		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":          err,
			"payment.id":     id,
			"payment.status": status,
		}, "[PaymentService][UpdatePaymentStatus] Failed to get payment by ID")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	payment.Status = status

	err = s.repo.UpdatePayment(ctx, tx, payment)
	if err != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":          err,
			"payment.id":     id,
			"payment.status": status,
		}, "[PaymentService][UpdatePaymentStatus] Failed to update payment")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	if status == enum.PaymentStatusSuccess {
		// Triggers
		for _, p := range payload {
			switch p.Type {
			case enum.PaymentTypeCourse:
				// TODO: trigger course enrollment
				fmt.Println("Trigger course enrollment")
			case enum.PaymentTypeMentor:
				// TODO: trigger mentorship booking
				fmt.Println("Trigger mentorship booking")
			case enum.PaymentTypeChallenge:
				// TODO: trigger challenge enrollment
				fmt.Println("Trigger challenge enrollment")
			}
		}
	}

	log.Info(map[string]interface{}{
		"payment.id":     id,
		"payment.status": status,
	}, "[PaymentService][UpdatePaymentStatus] Payment status updated")

	if err2 := tx.Commit(); err2 != nil {
		traceID := log.ErrorWithTraceID(map[string]interface{}{
			"error":          err2,
			"payment.id":     id,
			"payment.status": status,
		}, "[PaymentService][UpdatePaymentStatus] Failed to commit transaction")
		return errorpkg.ErrInternalServer.Build().WithTraceID(traceID)
	}

	return nil
}
