package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/usecase"
)

// UserCLI は CLI からの入力を usecase.RegisterUser に橋渡しする。
//
// handler.UserHandler と入出力チャネルが違うだけで責務は同じ。
// ただし注目したいのは Register() の switch:
// handler/user.go の Register() 内の switch とエラー分類の構造が
// 完全に同じで、出力先（http.Error vs Fprintf+ExitCode）だけが違う。
// ＝ レイヤードでは「エラー → 出力フォーマット」の翻訳ロジックが
// 入口の数だけ複製される。
type UserCLI struct {
	registerUser *usecase.RegisterUser
}

func NewUserCLI(registerUser *usecase.RegisterUser) *UserCLI {
	return &UserCLI{registerUser: registerUser}
}

// Register は登録処理を実行し、終了コードを返す。
func (c *UserCLI) Register(ctx context.Context, email, password string, stdout, stderr io.Writer) int {
	out, err := c.registerUser.Execute(ctx, usecase.RegisterUserInput{
		Email:    email,
		Password: password,
	})
	if err != nil {
		// ↓ handler/user.go:42-50 の switch と構造がそっくり同じ。
		//   出口が変わるたびに同じ分類ロジックを書き直す羽目になるのが
		//   Layered の弱点。
		switch {
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrPasswordTooShort):
			fmt.Fprintf(stderr, "invalid input: %v\n", err)
			return 2
		case errors.Is(err, user.ErrEmailAlreadyExists):
			fmt.Fprintf(stderr, "conflict: %v\n", err)
			return 3
		default:
			fmt.Fprintf(stderr, "internal error: %v\n", err)
			return 1
		}
	}

	fmt.Fprintf(stdout, "created user id=%d email=%s\n", out.ID, out.Email)
	return 0
}
