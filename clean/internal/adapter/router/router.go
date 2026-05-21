package router

import (
	"net/http"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/controller"
)

func New(userController *controller.UserController) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /users", userController.Register)

	return mux
}
