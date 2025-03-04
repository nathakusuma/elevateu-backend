package ctxkey

type contextKey string

const (
	UserID                contextKey = "user.id"
	UserRole              contextKey = "user.role"
	IsSubscribedBoost     contextKey = "user.is_subscribed_boost"
	IsSubscribedChallenge contextKey = "user.is_subscribed_challenge"
)

func (c contextKey) String() string {
	return string(c)
}
