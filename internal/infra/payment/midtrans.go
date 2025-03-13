package payment

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"

	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/domain/errorpkg"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
)

type midtransPayment struct{}

func NewMidtrans() IPaymentGateway {
	midtrans.ServerKey = env.GetEnv().MidtransServerKey
	midtrans.Environment = env.GetEnv().MidtransEnvironment
	return &midtransPayment{}
}

func (p *midtransPayment) CreateTransaction(id string, amount int) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  id,
			GrossAmt: int64(amount),
		},
		Expiry: &snap.ExpiryDetails{
			Unit:     "hour",
			Duration: 1,
		},
	}

	resp, err := snap.CreateTransaction(req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (p *midtransPayment) ProcessNotification(
	notificationPayload map[string]interface{}) (enum.PaymentStatus, string, error) {
	// 3. Get order-id from payload
	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		// do something when key `order_id` not found
		return "", "", errorpkg.ErrValidation.WithDetail("order_id not found")
	}

	statusCode, exists := notificationPayload["status_code"].(string)
	if !exists {
		return "", "", errorpkg.ErrValidation.WithDetail("status_code not found")
	}

	grossAmount, exists := notificationPayload["gross_amount"].(string)
	if !exists {
		return "", "", errorpkg.ErrValidation.WithDetail("gross_amount not found")
	}

	signatureKey, exists := notificationPayload["signature_key"].(string)
	if !exists {
		return "", "", errorpkg.ErrValidation.WithDetail("signature_key not found")
	}

	if ok := p.verifySignature(orderId, statusCode, grossAmount, signatureKey); !ok {
		return "", "", errorpkg.ErrValidation.WithDetail("signature_key not valid")
	}

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := coreapi.CheckTransaction(orderId)
	if e != nil {
		return "", "", errorpkg.ErrOKIgnore
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// set transaction status on your database to 'challenge'
					return enum.PaymentStatusChallenge, transactionStatusResp.PaymentType, nil
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// set transaction status on your database to 'success'
					return enum.PaymentStatusSuccess, transactionStatusResp.PaymentType, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// set transaction status on your databaase to 'success'
				return enum.PaymentStatusSuccess, transactionStatusResp.PaymentType, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// set transaction status on your databaase to 'failure'
				return enum.PaymentStatusFailure, transactionStatusResp.PaymentType, nil
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// set transaction status on your databaase to 'pending' / waiting payment
				return enum.PaymentStatusPending, transactionStatusResp.PaymentType, nil
			}
		}
	}

	return "", "", errorpkg.ErrInternalServer
}

func (p *midtransPayment) verifySignature(orderID, statusCode, grossAmount, providedSignature string) bool {
	serverKey := midtrans.ServerKey

	signatureString := orderID + statusCode + grossAmount + serverKey

	hasher := sha512.New()
	hasher.Write([]byte(signatureString))
	calculatedSignature := hex.EncodeToString(hasher.Sum(nil))

	return calculatedSignature == providedSignature
}
