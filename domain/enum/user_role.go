package enum

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleMentor  UserRole = "mentor"
	UserRoleStudent UserRole = "student"
)

func (r UserRole) String() string {
	return string(r)
}
