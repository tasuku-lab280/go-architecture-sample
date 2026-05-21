package server

import (
	"log"
	"net/http"
	"os"
)

func ListenAndServe(handler http.Handler) error {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server listening on :%s", port)
	return http.ListenAndServe(":"+port, handler)
}
