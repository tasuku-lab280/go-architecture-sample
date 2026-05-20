package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/handler"
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

func newTestHandler() *handler.UserHandler {
	uc := usecase.NewRegisterUser(newInMemoryUserRepository())
	return handler.NewUserHandler(uc)
}

func postUsers(h *handler.UserHandler, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.Register(rec, req)
	return rec
}

func TestUserHandler_Register(t *testing.T) {
	t.Run("正常系: 201でidとemailを返す", func(t *testing.T) {
		h := newTestHandler()

		rec := postUsers(h, map[string]string{
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
	})

	t.Run("異常系: 不正なJSONは400", func(t *testing.T) {
		h := newTestHandler()
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("not json")))
		rec := httptest.NewRecorder()
		h.Register(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status: got %d want 400", rec.Code)
		}
	})

	t.Run("異常系: 不正なメールアドレスは400", func(t *testing.T) {
		h := newTestHandler()
		rec := postUsers(h, map[string]string{
			"email":    "invalid",
			"password": "password123",
		})
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status: got %d want 400", rec.Code)
		}
	})

	t.Run("異常系: パスワードが短いと400", func(t *testing.T) {
		h := newTestHandler()
		rec := postUsers(h, map[string]string{
			"email":    "short@example.com",
			"password": "short",
		})
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status: got %d want 400", rec.Code)
		}
	})

	t.Run("異常系: 重複登録は409", func(t *testing.T) {
		h := newTestHandler()
		body := map[string]string{
			"email":    "dup@example.com",
			"password": "password123",
		}
		if rec := postUsers(h, body); rec.Code != http.StatusCreated {
			t.Fatalf("setup: got %d", rec.Code)
		}

		rec := postUsers(h, body)
		if rec.Code != http.StatusConflict {
			t.Errorf("status: got %d want 409", rec.Code)
		}
	})
}
