package usecase

import (
	"context"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
)

type RegisterUserInput struct {
	Email    string
	Password string
}

type RegisterUserOutput struct {
	ID    int64
	Email string
}

type RegisterUser struct {
	repo user.Repository
}

func NewRegisterUser(repo user.Repository) *RegisterUser {
	return &RegisterUser{repo: repo}
}

func (uc *RegisterUser) Execute(ctx context.Context, in RegisterUserInput) (*RegisterUserOutput, error) {
	u, err := user.NewUser(in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	exists, err := uc.repo.ExistsByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, user.ErrEmailAlreadyExists
	}

	if err := uc.repo.Save(ctx, u); err != nil {
		return nil, err
	}

	return &RegisterUserOutput{
		ID:    u.ID,
		Email: u.Email.String(),
	}, nil
}
