package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/infrastructure/database"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/usecase"
)

// Layered 版 CLI。
//
// 注目ポイント: ユースケース呼び出しは Handler と同じだが、
// 「エラー → 終了コード」「成功 → 標準出力フォーマット」の翻訳ロジックを
// ここで再実装している。handler/user.go の Register と見比べると、
// エラー判定の switch がほぼコピペ構造になっていることが分かる。
func main() {
	email := flag.String("email", "", "user email")
	password := flag.String("password", "", "user password")
	flag.Parse()

	if *email == "" || *password == "" {
		fmt.Fprintln(os.Stderr, "usage: cli -email <email> -password <password>")
		os.Exit(2)
	}

	db, err := database.NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect db: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := database.NewUserRepository(db)
	registerUser := usecase.NewRegisterUser(userRepo)

	out, err := registerUser.Execute(context.Background(), usecase.RegisterUserInput{
		Email:    *email,
		Password: *password,
	})
	if err != nil {
		// ↓ HTTP Handler の switch とそっくり同じ判定を、出口が違うだけで書き直している
		switch {
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrPasswordTooShort):
			fmt.Fprintf(os.Stderr, "invalid input: %v\n", err)
			os.Exit(2)
		case errors.Is(err, user.ErrEmailAlreadyExists):
			fmt.Fprintf(os.Stderr, "conflict: %v\n", err)
			os.Exit(3)
		default:
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("created user id=%d email=%s\n", out.ID, out.Email)
}
