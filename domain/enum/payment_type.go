package enum

type PaymentType string

const (
	PaymentTypeCourse    PaymentType = "course"
	PaymentTypeMentor    PaymentType = "mentor"
	PaymentTypeChallenge PaymentType = "challenge"
)
