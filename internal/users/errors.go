package users

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidRole  = errors.New("role is not valid")
)
