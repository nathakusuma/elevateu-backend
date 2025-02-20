package enum

type UserRole string

const (
	RoleAdmin   UserRole = "admin"
	RoleMentor  UserRole = "mentor"
	RoleStudent UserRole = "student"
)

func (r UserRole) String() string {
	return string(r)
}
