package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/handler"
)

func main() {
	mux := handler.NewMux()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
