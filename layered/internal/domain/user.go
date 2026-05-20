package domain

import (
	"context"
	"errors"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPasswordTooShort   = errors.New("password too short")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Email string

var emailRegexp = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func NewEmail(s string) (Email, error) {
	if !emailRegexp.MatchString(s) {
		return "", ErrInvalidEmail
	}
	return Email(s), nil
}

func (e Email) String() string { return string(e) }

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

type UserRepository interface {
	Save(ctx context.Context, u *User) error
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
