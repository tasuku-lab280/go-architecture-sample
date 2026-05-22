package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/gateway"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/presenter"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/infrastructure/database"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/interactor"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
)

// Clean 版 CLI。
//
// 注目ポイント:
//   - interactor.NewRegisterUserInteractor は HTTP 版 main.go と完全に同じ呼び出し
//   - 違うのは Presenter だけ（NewCLIRegisterUserPresenter）
//   - エラー → 終了コードの switch は Presenter 内に閉じている（このファイルには無い）
//
// 「ユースケースは触らない」「翻訳器（Presenter）を差し替えるだけ」が体感できる。
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

	userRepo := gateway.NewUserRepository(db)

	p := presenter.NewCLIRegisterUserPresenter()
	uc := interactor.NewRegisterUserInteractor(userRepo, p)

	_ = uc.Handle(context.Background(), input.RegisterUserInputData{
		Email:    *email,
		Password: *password,
	})

	p.WriteTo(os.Stdout, os.Stderr)
	os.Exit(p.ViewModel.ExitCode)
}
