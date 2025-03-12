package payment

import "github.com/nathakusuma/elevateu-backend/domain/enum"

type IPaymentGateway interface {
	CreateTransaction(id string, amount int) (string, error)
	ProcessNotification(notificationPayload map[string]interface{}) (enum.PaymentStatus, string, error)
}
