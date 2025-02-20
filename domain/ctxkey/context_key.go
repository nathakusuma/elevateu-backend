package ctxkey

type contextKey string

const (
	UserID   contextKey = "user.id"
	UserRole contextKey = "user.role"
)

func (c contextKey) String() string {
	return string(c)
}
