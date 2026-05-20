package interactor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/interactor"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/input"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/usecase/port/output"
)

type inMemoryUserRepository struct {
	users  map[user.Email]*user.User
	nextID int64
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{users: map[user.Email]*user.User{}}
}

func (r *inMemoryUserRepository) Save(_ context.Context, u *user.User) error {
	r.nextID++
	u.ID = r.nextID
	r.users[u.Email] = u
	return nil
}

func (r *inMemoryUserRepository) ExistsByEmail(_ context.Context, email user.Email) (bool, error) {
	_, ok := r.users[email]
	return ok, nil
}

type spyPresenter struct {
	out output.RegisterUserOutputData
	err error
}

func (p *spyPresenter) Present(out output.RegisterUserOutputData) { p.out = out }
func (p *spyPresenter) PresentError(err error)                    { p.err = err }

func TestRegisterUserInteractor_Handle(t *testing.T) {
	t.Run("正常系: Presenter.Present に出力データが渡される", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		p := &spyPresenter{}
		uc := interactor.NewRegisterUserInteractor(repo, p)

		err := uc.Handle(context.Background(), input.RegisterUserInputData{
			Email:    "new@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.out.ID == 0 {
			t.Error("ID should be assigned")
		}
		if p.out.Email != "new@example.com" {
			t.Errorf("email: got %q", p.out.Email)
		}
		if p.err != nil {
			t.Errorf("PresentError should not be called: %v", p.err)
		}
	})

	t.Run("異常系: 既存メールアドレスは Presenter.PresentError に ErrEmailAlreadyExists", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		uc := interactor.NewRegisterUserInteractor(repo, &spyPresenter{})
		in := input.RegisterUserInputData{Email: "dup@example.com", Password: "password123"}
		if err := uc.Handle(context.Background(), in); err != nil {
			t.Fatalf("setup: %v", err)
		}

		p := &spyPresenter{}
		uc2 := interactor.NewRegisterUserInteractor(repo, p)
		err := uc2.Handle(context.Background(), in)
		if !errors.Is(err, user.ErrEmailAlreadyExists) {
			t.Errorf("error: got %v want ErrEmailAlreadyExists", err)
		}
		if !errors.Is(p.err, user.ErrEmailAlreadyExists) {
			t.Errorf("PresentError: got %v want ErrEmailAlreadyExists", p.err)
		}
	})

	t.Run("異常系: 不正なメールアドレスはドメイン層のエラーが Presenter に伝搬", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		p := &spyPresenter{}
		uc := interactor.NewRegisterUserInteractor(repo, p)

		err := uc.Handle(context.Background(), input.RegisterUserInputData{
			Email:    "invalid",
			Password: "password123",
		})
		if !errors.Is(err, user.ErrInvalidEmail) {
			t.Errorf("error: got %v want ErrInvalidEmail", err)
		}
		if !errors.Is(p.err, user.ErrInvalidEmail) {
			t.Errorf("PresentError: got %v want ErrInvalidEmail", p.err)
		}
	})
}
