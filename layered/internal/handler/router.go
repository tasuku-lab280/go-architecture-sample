package handler

import "net/http"

func NewMux(userHandler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", userHandler.Register)

	return mux
}
