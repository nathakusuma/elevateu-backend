package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

type midtransService struct{}

func NewMidtransService() contract.IPaymentGateway {
	return &midtransService{}
}

func (s *midtransService) CreateTransaction(id string, amount int) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  id,
			GrossAmt: int64(amount),
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "minute",
			Duration: 15,
		},
	}

	resp, err := snap.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

// ProcessNotification is called as a callback from midtrans notification
func (s *midtransService) ProcessNotification(ctx context.Context, notificationPayload map[string]any,
	statusUpdateCallback func(context.Context, uuid.UUID, enum.PaymentStatus) error) error {
	// 1. Get order-id from payload
	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		return errorpkg.ErrValidation
	}

	// 2. Check transaction to Midtrans with param orderID
	transactionStatusResp, e := coreapi.CheckTransaction(orderID)
	if e != nil || transactionStatusResp == nil {
		return errorpkg.ErrInternalServer
	}

	// 3. Parse transaction ID
	transactionID, err := uuid.Parse(orderID)
	if err != nil {
		return errorpkg.ErrValidation
	}

	// 4. Update status based on transaction response
	return s.handleTransactionStatus(ctx, transactionID, transactionStatusResp, statusUpdateCallback, notificationPayload)
}

func (s *midtransService) handleTransactionStatus(
	ctx context.Context,
	transactionID uuid.UUID,
	transactionStatusResp *coreapi.TransactionStatusResponse,
	statusUpdateCallback func(context.Context, uuid.UUID, enum.PaymentStatus) error,
	notificationPayload map[string]any,
) error {
	transStatus := transactionStatusResp.TransactionStatus
	fraudStatus := transactionStatusResp.FraudStatus

	// Handle special case for "capture" status
	if transStatus == "capture" {
		return s.handleCaptureStatus(ctx, transactionID, fraudStatus, statusUpdateCallback, notificationPayload)
	}

	// Map transaction statuses to payment statuses
	statusMap := map[string]enum.PaymentStatus{
		"settlement": enum.PaymentStatusSuccess,
		"cancel":     enum.PaymentStatusFailure,
		"expire":     enum.PaymentStatusFailure,
		"pending":    enum.PaymentStatusPending,
	}

	// Get the corresponding payment status
	paymentStatus, exists := statusMap[transStatus]
	if !exists {
		// Handle "deny" status or other unhandled statuses
		if transStatus == "deny" {
			log.Info(map[string]interface{}{
				"notification_payload": notificationPayload,
			}, "[MidtransService][ProcessNotification] Payment status denied")
			return nil
		}
		// Unhandled status
		return errorpkg.ErrInternalServer
	}

	// Update the status
	return statusUpdateCallback(ctx, transactionID, paymentStatus)
}

// handleCaptureStatus processes the "capture" transaction status specifically
func (s *midtransService) handleCaptureStatus(
	ctx context.Context,
	transactionID uuid.UUID,
	fraudStatus string,
	statusUpdateCallback func(context.Context, uuid.UUID, enum.PaymentStatus) error,
	notificationPayload map[string]any,
) error {
	if fraudStatus == "challenge" {
		// Set transaction status to 'challenge'
		if err := statusUpdateCallback(ctx, transactionID, enum.PaymentStatusChallenge); err != nil {
			return err
		}

		log.Warn(map[string]interface{}{
			"notification_payload": notificationPayload,
		}, "[MidtransService][ProcessNotification] Payment status challenged. "+
			"Please take action on your Merchant Administration Portal")
		return nil
	} else if fraudStatus == "accept" {
		// Set transaction status to 'success'
		return statusUpdateCallback(ctx, transactionID, enum.PaymentStatusSuccess)
	}

	// Unhandled fraud status
	return errorpkg.ErrInternalServer
}
