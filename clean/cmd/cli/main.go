package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/controller"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/gateway"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/infrastructure/database"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/interactor"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

// Clean 版 CLI のエントリポイント。
//
// main の責務は DI 組み立てのみ:
//   - DB / Gateway / Interactor ファクトリの結線
//   - flag の解析（ここは外界との境界なので main に置く）
//   - CLIUserController に処理を委譲
//
// HTTP 版 cmd/api/main.go と比較すると、registerUserFactory までは
// 完全に同じ。違いは「どの Controller に渡すか」だけ。
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

	registerUserFactory := func(p output.RegisterUserPresenter) input.RegisterUserInputPort {
		return interactor.NewRegisterUserInteractor(userRepo, p)
	}

	cliController := controller.NewCLIUserController(registerUserFactory)
	code := cliController.Register(context.Background(), *email, *password, os.Stdout, os.Stderr)
	os.Exit(code)
}
