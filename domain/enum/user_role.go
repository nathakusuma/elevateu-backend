package enum

type UserRole string

const (
	RoleUser UserRole = "user"
)

func (r UserRole) String() string {
	return string(r)
}
