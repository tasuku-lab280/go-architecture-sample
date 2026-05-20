package user

import "regexp"

type Email string

var emailRegexp = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func NewEmail(s string) (Email, error) {
	if !emailRegexp.MatchString(s) {
		return "", ErrInvalidEmail
	}
	return Email(s), nil
}

func (e Email) String() string { return string(e) }
