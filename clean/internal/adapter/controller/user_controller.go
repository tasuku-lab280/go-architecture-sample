package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/presenter"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

// RegisterUserInteractorFactory は、リクエストごとに生成した Presenter を
// Interactor に束ねるためのファクトリ。
// 具象 Interactor を Controller が知らずに済ませるための間接層。
type RegisterUserInteractorFactory func(p output.RegisterUserPresenter) input.RegisterUserInputPort

type UserController struct {
	newRegisterUser RegisterUserInteractorFactory
}

func NewUserController(newRegisterUser RegisterUserInteractorFactory) *UserController {
	return &UserController{newRegisterUser: newRegisterUser}
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p := presenter.NewRegisterUserPresenter()
	uc := c.newRegisterUser(p)
	_ = uc.Handle(r.Context(), input.RegisterUserInputData{
		Email:    req.Email,
		Password: req.Password,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(p.ViewModel.StatusCode)
	_ = json.NewEncoder(w).Encode(p.ViewModel.Body)
}
