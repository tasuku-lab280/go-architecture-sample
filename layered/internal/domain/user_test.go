package domain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"正常系", "test@example.com", nil},
		{"@がない", "testexample.com", domain.ErrInvalidEmail},
		{"ドメイン部がない", "test@", domain.ErrInvalidEmail},
		{"TLDがない", "test@example", domain.ErrInvalidEmail},
		{"空文字", "", domain.ErrInvalidEmail},
		{"スペース混入", "test @example.com", domain.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.NewEmail(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error: got %v want %v", err, tt.wantErr)
			}
			if err == nil && got.String() != tt.input {
				t.Errorf("string: got %q want %q", got.String(), tt.input)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	t.Run("正常系: メールアドレスとパスワードが有効", func(t *testing.T) {
		u, err := domain.NewUser("test@example.com", "password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.Email.String() != "test@example.com" {
			t.Errorf("email: got %q", u.Email)
		}
		if u.PasswordHash == "password123" {
			t.Error("password must not be stored as plaintext")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("password123")); err != nil {
			t.Errorf("password hash should match original: %v", err)
		}
	})

	t.Run("異常系: 不正なメールアドレス", func(t *testing.T) {
		_, err := domain.NewUser("invalid", "password123")
		if !errors.Is(err, domain.ErrInvalidEmail) {
			t.Errorf("error: got %v want ErrInvalidEmail", err)
		}
	})

	t.Run("異常系: パスワードが7文字以下なら拒否", func(t *testing.T) {
		_, err := domain.NewUser("test@example.com", strings.Repeat("a", 7))
		if !errors.Is(err, domain.ErrPasswordTooShort) {
			t.Errorf("error: got %v want ErrPasswordTooShort", err)
		}
	})

	t.Run("境界値: パスワード8文字なら通る", func(t *testing.T) {
		_, err := domain.NewUser("test@example.com", strings.Repeat("a", 8))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
