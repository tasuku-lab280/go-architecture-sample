package interactor

import (
	"context"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

type RegisterUserInteractor struct {
	repo      output.UserRepository
	presenter output.RegisterUserPresenter
}

func NewRegisterUserInteractor(repo output.UserRepository, presenter output.RegisterUserPresenter) input.RegisterUserInputPort {
	return &RegisterUserInteractor{repo: repo, presenter: presenter}
}

func (i *RegisterUserInteractor) Handle(ctx context.Context, in input.RegisterUserInputData) error {
	u, err := user.NewUser(in.Email, in.Password)
	if err != nil {
		i.presenter.PresentError(err)
		return err
	}

	exists, err := i.repo.ExistsByEmail(ctx, u.Email)
	if err != nil {
		i.presenter.PresentError(err)
		return err
	}
	if exists {
		i.presenter.PresentError(user.ErrEmailAlreadyExists)
		return user.ErrEmailAlreadyExists
	}

	if err := i.repo.Save(ctx, u); err != nil {
		i.presenter.PresentError(err)
		return err
	}

	i.presenter.Present(output.RegisterUserOutputData{
		ID:    u.ID,
		Email: u.Email.String(),
	})
	return nil
}
