package users

type Role int

const (
	RoleAdmin Role = 0
	RoleUser  Role = 10
)

type User struct {
	Email string `json:"email"`
	Roles []Role `json:"roles"`
}

func (u *User) HasRole(r Role) bool {
	for _, role := range u.Roles {
		if int(role) >= int(r) {
			return true
		}
	}
	return false
}
