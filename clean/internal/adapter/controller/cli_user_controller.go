package controller

import (
	"context"
	"io"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/presenter"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
)

// CLIUserController は CLI からの入力をユースケース入力ポートに翻訳する。
//
// HTTP 用 UserController と入出力が違うだけで、責務は同じ:
//   - 外界の入力（flag 文字列）を input.RegisterUserInputData に変換
//   - リクエストごとに Presenter を生成し Interactor に注入
//   - Interactor を呼び、ViewModel を外界（stdout / stderr / 終了コード）に書き出す
//
// HTTP 版と同じ RegisterUserInteractorFactory 型を再利用できるのがポイント。
// ファクトリは Presenter を受け取って入力ポートを返すだけで、
// その Presenter が HTTP 用か CLI 用かは関知しない。
type CLIUserController struct {
	newRegisterUser RegisterUserInteractorFactory
}

func NewCLIUserController(newRegisterUser RegisterUserInteractorFactory) *CLIUserController {
	return &CLIUserController{newRegisterUser: newRegisterUser}
}

// Register は登録処理を実行し、終了コードを返す。
// stdout / stderr は呼び出し側（main）が注入する。テストで bytes.Buffer に差し替え可能。
func (c *CLIUserController) Register(ctx context.Context, email, password string, stdout, stderr io.Writer) int {
	p := presenter.NewCLIRegisterUserPresenter()
	uc := c.newRegisterUser(p)
	_ = uc.Handle(ctx, input.RegisterUserInputData{
		Email:    email,
		Password: password,
	})

	p.WriteTo(stdout, stderr)
	return p.ViewModel.ExitCode
}
