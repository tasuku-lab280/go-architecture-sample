package user_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
)

func TestNewPassword(t *testing.T) {
	t.Run("正常系: 平文はハッシュ化されて保持される", func(t *testing.T) {
		p, err := user.NewPassword("password123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Hash() == "password123" {
			t.Error("password must not be stored as plaintext")
		}
	})

	t.Run("異常系: 7文字以下なら拒否", func(t *testing.T) {
		_, err := user.NewPassword(strings.Repeat("a", 7))
		if !errors.Is(err, user.ErrPasswordTooShort) {
			t.Errorf("error: got %v want ErrPasswordTooShort", err)
		}
	})

	t.Run("境界値: 8文字なら通る", func(t *testing.T) {
		_, err := user.NewPassword(strings.Repeat("a", 8))
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestPassword_Verify(t *testing.T) {
	t.Run("同じ平文と一致する", func(t *testing.T) {
		p, _ := user.NewPassword("password123")
		if !p.Verify("password123") {
			t.Error("verify should be true for the same plaintext")
		}
	})

	t.Run("異なる平文とは一致しない", func(t *testing.T) {
		p, _ := user.NewPassword("password123")
		if p.Verify("wrongpassword") {
			t.Error("verify should be false for a different plaintext")
		}
	})
}
