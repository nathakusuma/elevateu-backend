package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/cache"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/internal/infra/payment"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
	"github.com/nathakusuma/elevateu-backend/pkg/uuidpkg"
)

type paymentService struct {
	repo           contract.IPaymentRepository
	mentoringSvc   contract.IMentoringService
	userSvc        contract.IUserService
	cache          cache.ICache
	paymentGateway payment.IPaymentGateway
	txManager      database.ITransactionManager
	uuid           uuidpkg.IUUID
}

func NewPaymentService(
	repo contract.IPaymentRepository,
	mentoringSvc contract.IMentoringService,
	userSvc contract.IUserService,
	cache cache.ICache,
	paymentGateway payment.IPaymentGateway,
	txManager database.ITransactionManager,
	uuid uuidpkg.IUUID,
) contract.IPaymentService {
	return &paymentService{
		repo:           repo,
		mentoringSvc:   mentoringSvc,
		userSvc:        userSvc,
		cache:          cache,
		paymentGateway: paymentGateway,
		txManager:      txManager,
		uuid:           uuid,
	}
}

func (s *paymentService) createPayment(ctx context.Context, req dto.CreatePaymentRequest) (string, error) {
	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to begin transaction")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	paymentID, err := s.uuid.NewV7()
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to generate payment ID")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	token, err := s.paymentGateway.CreateTransaction(paymentID.String(), req.Amount)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to create transaction in payment gateway")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	paymentEntity := &entity.Payment{
		ID:        paymentID,
		UserID:    req.UserID,
		Token:     token,
		Amount:    req.Amount,
		Title:     req.Title,
		Detail:    req.Detail,
		Status:    enum.PaymentStatusPending,
		ExpiredAt: time.Now().Add(1 * time.Hour),
	}

	if err2 := s.repo.CreatePayment(ctx, tx, paymentEntity); err2 != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err2,
			"request": req,
		}, "Failed to create payment")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	payloadJSON, err := sonic.Marshal(req.Payload)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to marshal payment payload")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err := s.cache.Set(ctx, "payment:"+paymentEntity.ID.String(), string(payloadJSON), 1*time.Hour); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to set payment payload in cache")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err := tx.Commit(); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":   err,
			"request": req,
		}, "Failed to commit transaction")
		return "", errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"payment.id":    paymentID,
		"payment.token": token,
		"request":       req,
	}, "Payment created")

	return token, nil
}

func (s *paymentService) UpdatePaymentStatus(ctx context.Context, id uuid.UUID, status enum.PaymentStatus,
	method string) error {
	tx, err := s.txManager.BeginTx(ctx)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":          err,
			"payment.id":     id,
			"payment.status": status,
		}, "Failed to begin transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}
	defer tx.Rollback()

	paymentEntity, err := s.repo.GetPaymentByID(ctx, tx, id)
	if err != nil {
		if strings.HasPrefix(err.Error(), "payment payload not found") {
			return errorpkg.ErrNotFound()
		}

		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":          err,
			"payment.id":     id,
			"payment.status": status,
		}, "Failed to get payment by ID")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	paymentEntity.Status = status
	paymentEntity.Method = method

	err = s.repo.UpdatePayment(ctx, tx, paymentEntity)
	if err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":          err,
			"payment.id":     id,
			"payment.status": status,
		}, "Failed to update payment")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if status == enum.PaymentStatusSuccess {
		var payloadJSON string
		if err = s.cache.Get(ctx, "payment:"+id.String(), &payloadJSON); err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":          err,
				"payment.id":     id,
				"payment.status": status,
			}, "Failed to get payment payload")
			return errorpkg.ErrInternalServer().WithTraceID(traceID)
		}

		var payload entity.PaymentPayload
		if err = sonic.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
				"error":          err,
				"payment.id":     id,
				"payment.status": status,
			}, "Failed to unmarshal payment payload")
			return errorpkg.ErrInternalServer().WithTraceID(traceID)
		}

		// Triggers
		switch payload.Type {
		case enum.PaymentTypeBoost:
			if err := s.triggerSkillBoost(ctx, tx, payload, paymentEntity); err != nil {
				return err
			}
		case enum.PaymentTypeChallenge:
			if err := s.triggerSkillChallenge(ctx, tx, payload, paymentEntity); err != nil {
				return err
			}
		case enum.PaymentTypeGuidance:
			if err := s.triggerSkillGuidance(ctx, tx, payload, paymentEntity); err != nil {
				return err
			}
		}
	}

	if err2 := tx.Commit(); err2 != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":          err2,
			"payment.id":     id,
			"payment.status": status,
		}, "Failed to commit transaction")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"payment.id":     id,
		"payment.status": status,
	}, "Payment status updated")

	return nil
}

func (s *paymentService) ProcessNotification(ctx context.Context, notificationPayload map[string]any) error {
	status, method, err := s.paymentGateway.ProcessNotification(notificationPayload)
	if err != nil {
		return err
	}

	orderIDStr, exists := notificationPayload["order_id"].(string)
	if !exists {
		return errorpkg.ErrValidation().WithDetail("order_id not found")
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return errorpkg.ErrValidation().WithDetail("order_id not valid")
	}

	log.Info(ctx, map[string]interface{}{
		"payload": notificationPayload,
		"status":  status,
		"method":  method,
	}, "incoming payment notification")

	return s.UpdatePaymentStatus(ctx, orderID, status, method)
}

func (s *paymentService) PaySkillBoost(ctx context.Context, studentID uuid.UUID) (string, error) {
	user, err := s.userSvc.GetUserByID(ctx, studentID, false)
	if err != nil {
		return "", err
	}

	badge := user.Student.Badge

	price := 120000
	switch badge {
	case enum.BadgeBronze:
		// 10% discount
		price = price - (price * 10 / 100)
	case enum.BadgeSilver:
		// 20% discount
		price = price - (price * 20 / 100)
	case enum.BadgeGold:
		// 50% discount
		price = price - (price * 50 / 100)
	}

	detail := "Skill Boost Subscription for 30 days"
	return s.createPayment(ctx, dto.CreatePaymentRequest{
		UserID: studentID,
		Amount: price,
		Title:  "Skill Boost Subscription",
		Detail: &detail,
		Payload: entity.PaymentPayload{
			Type:      enum.PaymentTypeBoost,
			StudentID: studentID,
		},
	})
}

func (s *paymentService) PaySkillChallenge(ctx context.Context, studentID uuid.UUID) (string, error) {
	user, err := s.userSvc.GetUserByID(ctx, studentID, false)
	if err != nil {
		return "", err
	}

	badge := user.Student.Badge

	price := 120000
	switch badge {
	case enum.BadgeBronze:
		// 10% discount
		price = price - (price * 10 / 100)
	case enum.BadgeSilver:
		// 20% discount
		price = price - (price * 20 / 100)
	case enum.BadgeGold:
		// 50% discount
		price = price - (price * 50 / 100)
	}

	detail := "Skill Challenge Subscription for 30 days"
	return s.createPayment(ctx, dto.CreatePaymentRequest{
		UserID: studentID,
		Amount: price,
		Title:  "Skill Challenge Subscription",
		Detail: &detail,
		Payload: entity.PaymentPayload{
			Type:      enum.PaymentTypeChallenge,
			StudentID: studentID,
		},
	})
}

func (s *paymentService) PaySkillGuidance(ctx context.Context, studentID, mentorID uuid.UUID) (string, error) {
	// check if mentor exists
	mentor, err := s.userSvc.GetUserByID(ctx, mentorID, false)
	if err != nil {
		return "", err
	}

	detail := fmt.Sprintf("Skill Guidance with %s for 24 hours", mentor.Name)
	return s.createPayment(ctx, dto.CreatePaymentRequest{
		UserID: studentID,
		Amount: mentor.Mentor.Price,
		Title:  "Skill Guidance Subscription",
		Detail: &detail,
		Payload: entity.PaymentPayload{
			Type:      enum.PaymentTypeGuidance,
			StudentID: studentID,
			MentorID:  mentorID,
		},
	})
}

func (s *paymentService) triggerSkillBoost(ctx context.Context, tx database.ITransaction,
	payload entity.PaymentPayload, payment *entity.Payment) error {
	student, err := s.userSvc.GetUserByID(ctx, payload.StudentID, false)
	if err != nil {
		return err
	}

	var subscribedUntil time.Time
	if student.Student.SubscribedBoostUntil == nil || student.Student.SubscribedBoostUntil.Before(time.Now()) {
		subscribedUntil = time.Now().Add(time.Hour * 24 * 30)
	} else {
		subscribedUntil = student.Student.SubscribedBoostUntil.Add(time.Hour * 24 * 30)
	}

	if err := s.repo.AddBoostSubscription(ctx, tx, payload.StudentID, subscribedUntil); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":          err,
			"payment.id":     payment.ID,
			"payment.status": payment.Status,
		}, "Failed to add boost subscription")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"student.id":       student.ID,
		"subscribed.until": subscribedUntil,
	}, "Skill Boost subscription added")

	return nil
}

func (s *paymentService) triggerSkillChallenge(ctx context.Context, tx database.ITransaction,
	payload entity.PaymentPayload, payment *entity.Payment) error {
	student, err := s.userSvc.GetUserByID(ctx, payload.StudentID, false)
	if err != nil {
		return err
	}

	var subscribedUntil time.Time
	if student.Student.SubscribedChallengeUntil == nil || student.Student.SubscribedChallengeUntil.Before(time.Now()) {
		subscribedUntil = time.Now().Add(time.Hour * 24 * 30)
	} else {
		subscribedUntil = student.Student.SubscribedChallengeUntil.Add(time.Hour * 24 * 30)
	}

	if err := s.repo.AddChallengeSubscription(ctx, tx, payload.StudentID, subscribedUntil); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"error":          err,
			"payment.id":     payment.ID,
			"payment.status": payment.Status,
		}, "Failed to add challenge subscription")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	log.Info(ctx, map[string]interface{}{
		"student.id":       student.ID,
		"subscribed.until": subscribedUntil,
	}, "Skill Challenge subscription added")

	return nil
}

func (s *paymentService) triggerSkillGuidance(ctx context.Context, tx database.ITransaction,
	payload entity.PaymentPayload, payment *entity.Payment) error {
	student, err := s.userSvc.GetUserByID(ctx, payload.StudentID, false)
	if err != nil {
		return err
	}

	mentor, err := s.userSvc.GetUserByID(ctx, payload.MentorID, false)
	if err != nil {
		return err
	}

	detail := fmt.Sprintf("Skill Guidance with %s", student.Name)
	mentorSalary := mentor.Mentor.Price - (mentor.Mentor.Price * 5 / 100) // Potongan 5% untuk ElevateU
	if err := s.repo.CreateMentorTransactionHistory(ctx, tx, &entity.MentorTransactionHistory{
		ID:       payment.ID,
		MentorID: payload.MentorID,
		Title:    "Pembayaran Mentor",
		Detail:   &detail,
		Amount:   mentorSalary,
	}); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"payment.id":     payment.ID,
			"payment.status": payment.Status,
		}, "Failed to create mentor transaction history")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if err := s.repo.AddMentorBalance(ctx, tx, payload.MentorID, mentorSalary); err != nil {
		traceID := log.ErrorWithTraceID(ctx, map[string]interface{}{
			"payment.id":     payment.ID,
			"payment.status": payment.Status,
		}, "Failed to add mentor balance")
		return errorpkg.ErrInternalServer().WithTraceID(traceID)
	}

	if _, err := s.mentoringSvc.CreateChat(ctx, payload.MentorID, payload.StudentID, false); err != nil {
		return err
	}

	log.Info(ctx, map[string]interface{}{
		"student.id": student.ID,
		"mentor.id":  mentor.ID,
	}, "Skill Guidance subscription added")

	return nil
}
