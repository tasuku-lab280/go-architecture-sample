package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/usecase"
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

func TestRegisterUser_Execute(t *testing.T) {
	t.Run("正常系: 新規ユーザーを登録できる", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		uc := usecase.NewRegisterUser(repo)

		out, err := uc.Execute(context.Background(), usecase.RegisterUserInput{
			Email:    "new@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if out.ID == 0 {
			t.Error("ID should be assigned")
		}
		if out.Email != "new@example.com" {
			t.Errorf("email: got %q", out.Email)
		}
	})

	t.Run("異常系: 既存メールアドレスは409相当のエラー", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		uc := usecase.NewRegisterUser(repo)

		in := usecase.RegisterUserInput{Email: "dup@example.com", Password: "password123"}
		if _, err := uc.Execute(context.Background(), in); err != nil {
			t.Fatalf("setup: %v", err)
		}
		_, err := uc.Execute(context.Background(), in)
		if !errors.Is(err, user.ErrEmailAlreadyExists) {
			t.Errorf("error: got %v want ErrEmailAlreadyExists", err)
		}
	})

	t.Run("異常系: 不正なメールアドレスはドメイン層のエラーを伝搬", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		uc := usecase.NewRegisterUser(repo)

		_, err := uc.Execute(context.Background(), usecase.RegisterUserInput{
			Email:    "invalid",
			Password: "password123",
		})
		if !errors.Is(err, user.ErrInvalidEmail) {
			t.Errorf("error: got %v want ErrInvalidEmail", err)
		}
	})
}
