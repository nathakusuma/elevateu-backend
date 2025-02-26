package enum

type PaymentStatus string

const (
	PaymentStatusSuccess   PaymentStatus = "success"
	PaymentStatusFailure   PaymentStatus = "failure"
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusChallenge PaymentStatus = "challenge"
)
