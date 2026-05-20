package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/handler"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/infrastructure"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/usecase"
)

func main() {
	db, err := infrastructure.NewDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	userRepo := infrastructure.NewUserRepository(db)
	registerUser := usecase.NewRegisterUser(userRepo)
	userHandler := handler.NewUserHandler(registerUser)

	mux := handler.NewMux(userHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
