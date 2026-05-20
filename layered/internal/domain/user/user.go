package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Email        Email
	PasswordHash string
	CreatedAt    time.Time
}

const passwordMinLength = 8

func NewUser(email, password string) (*User, error) {
	addr, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	if len(password) < passwordMinLength {
		return nil, ErrPasswordTooShort
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:        addr,
		PasswordHash: string(hash),
	}, nil
}
