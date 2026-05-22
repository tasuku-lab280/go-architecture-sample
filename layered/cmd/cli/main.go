package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/cli"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/infrastructure/database"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/usecase"
)

// Layered 版 CLI のエントリポイント。
//
// main の責務は DI 組み立てのみ。
// 「エラー → 終了コード」「成功 → 標準出力フォーマット」の翻訳は
// internal/cli/user_cli.go に閉じている。
//
// internal/cli/user_cli.go と internal/handler/user.go を見比べると、
// エラー分類の switch がほぼ同じ形でコピペ的に重複している。
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
	userCLI := cli.NewUserCLI(registerUser)

	code := userCLI.Register(context.Background(), *email, *password, os.Stdout, os.Stderr)
	os.Exit(code)
}
