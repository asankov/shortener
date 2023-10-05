package users

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Role represents a user role.
//
// Each role has a set of actions that can perform.
//
// The level of permissions is shown by the int value of the role.
// The lower the value the more permissions a role has.
type Role int

//go:generate stringer -type=Role -linecomment
const (
	// RoleAdmin is a role representing an admin.
	// Its value is 0, which means no role can have greater permissions that it.
	RoleAdmin Role = 0 // Admin
	// RoleUser is a role representing a user.
	RoleUser Role = 10 // User
)

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(b []byte) error {
	roleString := string(b)
	roleString = strings.TrimRight(roleString, `"`)
	roleString = strings.TrimLeft(roleString, `"`)
	role, err := RoleFrom(roleString)
	if err != nil {
		return err
	}
	*r = role

	return nil
}

type User struct {
	Email string `json:"email"`
	Roles []Role `json:"roles"`
}

func (u *User) HasRole(r Role) bool {
	for _, role := range u.Roles {
		if int(role) <= int(r) {
			return true
		}
	}
	return false
}

func RoleFrom(s string) (Role, error) {
	switch strings.ToLower(s) {
	case "admin":
		return RoleAdmin, nil
	case "user":
		return RoleUser, nil
	default:
		return -1, fmt.Errorf("%w: %v", ErrInvalidRole, s)
	}
}
