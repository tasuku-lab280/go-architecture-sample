package user

import "time"

type User struct {
	ID        int64
	Email     Email
	Password  Password
	CreatedAt time.Time
}

func NewUser(email, password string) (*User, error) {
	addr, err := NewEmail(email)
	if err != nil {
		return nil, err
	}
	pw, err := NewPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{
		Email:    addr,
		Password: pw,
	}, nil
}
