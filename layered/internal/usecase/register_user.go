package usecase

import (
	"context"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain"
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
	repo domain.UserRepository
}

func NewRegisterUser(repo domain.UserRepository) *RegisterUser {
	return &RegisterUser{repo: repo}
}

func (uc *RegisterUser) Execute(ctx context.Context, in RegisterUserInput) (*RegisterUserOutput, error) {
	user, err := domain.NewUser(in.Email, in.Password)
	if err != nil {
		return nil, err
	}

	exists, err := uc.repo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	if err := uc.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return &RegisterUserOutput{
		ID:    user.ID,
		Email: user.Email.String(),
	}, nil
}
