package presenter

import (
	"errors"
	"net/http"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

// RegisterUserViewModel は Controller が HTTP レスポンスとして書き出す中間表現。
type RegisterUserViewModel struct {
	StatusCode int
	Body       any
}

type successBody struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type errorBody struct {
	Error string `json:"error"`
}

// RegisterUserPresenter は output.RegisterUserPresenter を実装し、
// ユースケース出力を HTTP 表現（ViewModel）に変換する。
// リクエストごとに生成して Interactor に注入する想定。
type RegisterUserPresenter struct {
	ViewModel RegisterUserViewModel
}

func NewRegisterUserPresenter() *RegisterUserPresenter {
	return &RegisterUserPresenter{}
}

func (p *RegisterUserPresenter) Present(out output.RegisterUserOutputData) {
	p.ViewModel = RegisterUserViewModel{
		StatusCode: http.StatusCreated,
		Body:       successBody{ID: out.ID, Email: out.Email},
	}
}

func (p *RegisterUserPresenter) PresentError(err error) {
	var code int
	switch {
	case errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrPasswordTooShort):
		code = http.StatusBadRequest
	case errors.Is(err, user.ErrEmailAlreadyExists):
		code = http.StatusConflict
	default:
		code = http.StatusInternalServerError
	}
	p.ViewModel = RegisterUserViewModel{
		StatusCode: code,
		Body:       errorBody{Error: err.Error()},
	}
}
