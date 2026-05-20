package main

import (
	"log"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/controller"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/gateway"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/router"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/infrastructure/database"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/infrastructure/server"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/interactor"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

func main() {
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	userRepo := gateway.NewUserRepository(db)

	// 具象 Interactor を Controller から隠すためのファクトリ。
	// リクエストごとに渡される Presenter と Repository シングルトンを束ねる。
	registerUserFactory := func(p output.RegisterUserPresenter) input.RegisterUserInputPort {
		return interactor.NewRegisterUserInteractor(userRepo, p)
	}

	userController := controller.NewUserController(registerUserFactory)
	mux := router.New(userController)

	if err := server.ListenAndServe(mux); err != nil {
		log.Fatal(err)
	}
}
