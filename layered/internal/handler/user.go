package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/usecase"
)

type UserHandler struct {
	registerUser *usecase.RegisterUser
}

func NewUserHandler(registerUser *usecase.RegisterUser) *UserHandler {
	return &UserHandler{registerUser: registerUser}
}

type registerUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerUserResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	out, err := h.registerUser.Execute(r.Context(), usecase.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidEmail),
			errors.Is(err, domain.ErrPasswordTooShort):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, domain.ErrEmailAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(registerUserResponse{
		ID:    out.ID,
		Email: out.Email,
	})
}
