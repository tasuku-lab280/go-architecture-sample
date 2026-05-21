package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/controller"
	"github.com/kudoutasuku/go-architecture-sample/clean/internal/adapter/presenter"
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

func newTestController() *controller.UserController {
	repo := newInMemoryUserRepository()
	factory := func(p output.RegisterUserPresenter) input.RegisterUserInputPort {
		return interactor.NewRegisterUserInteractor(repo, p)
	}
	return controller.NewUserController(factory)
}

func newTestControllerWith(repo output.UserRepository) *controller.UserController {
	factory := func(p output.RegisterUserPresenter) input.RegisterUserInputPort {
		return interactor.NewRegisterUserInteractor(repo, p)
	}
	return controller.NewUserController(factory)
}

func postUsers(c *controller.UserController, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	c.Register(rec, req)
	return rec
}

func TestUserController_Register(t *testing.T) {
	t.Run("正常系: 201でidとemailを返す", func(t *testing.T) {
		c := newTestController()

		rec := postUsers(c, map[string]string{
			"email":    "ok@example.com",
			"password": "password123",
		})

		if rec.Code != http.StatusCreated {
			t.Fatalf("status: got %d want 201, body=%s", rec.Code, rec.Body.String())
		}
		var resp struct {
			ID    int64  `json:"id"`
			Email string `json:"email"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if resp.ID == 0 || resp.Email != "ok@example.com" {
			t.Errorf("response: %+v", resp)
		}

		// Presenter が ViewModel.Body の Content-Type を JSON にしている
		if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
			t.Errorf("content-type: got %q", ct)
		}
	})

	t.Run("異常系: 不正なJSONは400", func(t *testing.T) {
		c := newTestController()
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("not json")))
		rec := httptest.NewRecorder()
		c.Register(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status: got %d want 400", rec.Code)
		}
	})

	t.Run("異常系: 不正なメールアドレスは400", func(t *testing.T) {
		c := newTestController()
		rec := postUsers(c, map[string]string{
			"email":    "invalid",
			"password": "password123",
		})
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status: got %d want 400", rec.Code)
		}
	})

	t.Run("異常系: パスワードが短いと400", func(t *testing.T) {
		c := newTestController()
		rec := postUsers(c, map[string]string{
			"email":    "short@example.com",
			"password": "short",
		})
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status: got %d want 400", rec.Code)
		}
	})

	t.Run("異常系: 重複登録は409", func(t *testing.T) {
		repo := newInMemoryUserRepository()
		c := newTestControllerWith(repo)
		body := map[string]string{
			"email":    "dup@example.com",
			"password": "password123",
		}
		if rec := postUsers(c, body); rec.Code != http.StatusCreated {
			t.Fatalf("setup: got %d", rec.Code)
		}

		rec := postUsers(c, body)
		if rec.Code != http.StatusConflict {
			t.Errorf("status: got %d want 409", rec.Code)
		}
	})
}

// Presenter が main と同じ型で組まれていることを念のため検証する。
var _ output.RegisterUserPresenter = (*presenter.RegisterUserPresenter)(nil)
