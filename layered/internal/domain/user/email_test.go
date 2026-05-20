package user_test

import (
	"errors"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"正常系", "test@example.com", nil},
		{"@がない", "testexample.com", user.ErrInvalidEmail},
		{"ドメイン部がない", "test@", user.ErrInvalidEmail},
		{"TLDがない", "test@example", user.ErrInvalidEmail},
		{"空文字", "", user.ErrInvalidEmail},
		{"スペース混入", "test @example.com", user.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewEmail(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error: got %v want %v", err, tt.wantErr)
			}
			if err == nil && got.String() != tt.input {
				t.Errorf("string: got %q want %q", got.String(), tt.input)
			}
		})
	}
}
