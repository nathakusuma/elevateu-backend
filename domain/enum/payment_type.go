package enum

type PaymentType string

const (
	PaymentTypeBoost     PaymentType = "boost"
	PaymentTypeChallenge PaymentType = "challenge"
	PaymentTypeGuidance  PaymentType = "guidance"
)
