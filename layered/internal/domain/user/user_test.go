package user_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
)

func TestNewUser(t *testing.T) {
	t.Run("正常系: メールアドレスとパスワードが有効", func(t *testing.T) {
		u, err := user.NewUser("test@example.com", "password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.Email.String() != "test@example.com" {
			t.Errorf("email: got %q", u.Email)
		}
		if !u.Password.Verify("password123") {
			t.Error("password should verify against the original plaintext")
		}
	})

	t.Run("異常系: 不正なメールアドレス", func(t *testing.T) {
		_, err := user.NewUser("invalid", "password123")
		if !errors.Is(err, user.ErrInvalidEmail) {
			t.Errorf("error: got %v want ErrInvalidEmail", err)
		}
	})

	t.Run("異常系: パスワードが7文字以下なら拒否", func(t *testing.T) {
		_, err := user.NewUser("test@example.com", strings.Repeat("a", 7))
		if !errors.Is(err, user.ErrPasswordTooShort) {
			t.Errorf("error: got %v want ErrPasswordTooShort", err)
		}
	})
}
