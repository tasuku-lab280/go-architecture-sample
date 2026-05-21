package user

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPasswordTooShort   = errors.New("password too short")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
