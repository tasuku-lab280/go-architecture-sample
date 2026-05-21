package presenter_test

import (
	"net/http"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/presenter"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

func TestRegisterUserPresenter_Present(t *testing.T) {
	p := presenter.NewRegisterUserPresenter()
	p.Present(output.RegisterUserOutputData{ID: 1, Email: "ok@example.com"})

	if p.ViewModel.StatusCode != http.StatusCreated {
		t.Errorf("status: got %d want 201", p.ViewModel.StatusCode)
	}
	if p.ViewModel.Body == nil {
		t.Error("body should be set")
	}
}

func TestRegisterUserPresenter_PresentError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"不正メール → 400", user.ErrInvalidEmail, http.StatusBadRequest},
		{"短いパスワード → 400", user.ErrPasswordTooShort, http.StatusBadRequest},
		{"重複メール → 409", user.ErrEmailAlreadyExists, http.StatusConflict},
		{"その他エラー → 500", errInternal{}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := presenter.NewRegisterUserPresenter()
			p.PresentError(tt.err)
			if p.ViewModel.StatusCode != tt.want {
				t.Errorf("status: got %d want %d", p.ViewModel.StatusCode, tt.want)
			}
		})
	}
}

type errInternal struct{}

func (errInternal) Error() string { return "internal" }
