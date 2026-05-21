package user

import "golang.org/x/crypto/bcrypt"

const passwordMinLength = 8

type Password struct {
	hash string
}

func NewPassword(plain string) (Password, error) {
	if len(plain) < passwordMinLength {
		return Password{}, ErrPasswordTooShort
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}
	return Password{hash: string(h)}, nil
}

func (p Password) Hash() string { return p.hash }

func (p Password) Verify(plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plain)) == nil
}
