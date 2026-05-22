package presenter

import (
	"errors"
	"fmt"
	"io"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

// CLIRegisterUserViewModel は CLI が標準出力／終了コードとして書き出す中間表現。
type CLIRegisterUserViewModel struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// CLIRegisterUserPresenter は output.RegisterUserPresenter を実装し、
// ユースケース出力を CLI 表現（テキスト + 終了コード）に変換する。
//
// HTTP 用の RegisterUserPresenter とは「変換先」が違うだけで、
// Interactor から見ると同じ output.RegisterUserPresenter インタフェース。
// つまり Interactor は一切変更せず、main で渡す Presenter を差し替えるだけで CLI 化できる。
type CLIRegisterUserPresenter struct {
	ViewModel CLIRegisterUserViewModel
}

func NewCLIRegisterUserPresenter() *CLIRegisterUserPresenter {
	return &CLIRegisterUserPresenter{}
}

func (p *CLIRegisterUserPresenter) Present(out output.RegisterUserOutputData) {
	p.ViewModel = CLIRegisterUserViewModel{
		Stdout:   fmt.Sprintf("created user id=%d email=%s\n", out.ID, out.Email),
		ExitCode: 0,
	}
}

func (p *CLIRegisterUserPresenter) PresentError(err error) {
	var code int
	var label string
	switch {
	case errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrPasswordTooShort):
		code = 2
		label = "invalid input"
	case errors.Is(err, user.ErrEmailAlreadyExists):
		code = 3
		label = "conflict"
	default:
		code = 1
		label = "internal error"
	}
	p.ViewModel = CLIRegisterUserViewModel{
		Stderr:   fmt.Sprintf("%s: %v\n", label, err),
		ExitCode: code,
	}
}

// WriteTo は ViewModel を実際の出力先に流し込む。
// 標準出力／標準エラーへの書き込みは副作用なので、Presenter 本体ではなく
// この関数（CLI ドライバ側）から呼ぶ形にしている。
func (p *CLIRegisterUserPresenter) WriteTo(stdout, stderr io.Writer) {
	if p.ViewModel.Stdout != "" {
		_, _ = io.WriteString(stdout, p.ViewModel.Stdout)
	}
	if p.ViewModel.Stderr != "" {
		_, _ = io.WriteString(stderr, p.ViewModel.Stderr)
	}
}
