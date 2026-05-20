package database_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain/user"
	"github.com/kudoutasuku/go-architecture-sample/layered/internal/infrastructure/database"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	setDefaultEnv("DB_HOST", "127.0.0.1")
	setDefaultEnv("DB_PORT", "3306")
	setDefaultEnv("DB_USER", "app")
	setDefaultEnv("DB_PASSWORD", "password")
	setDefaultEnv("DB_NAME", "app_db")

	if db, err := database.NewDB(); err == nil {
		testDB = db
		defer testDB.Close()
		if err := ensureSchema(testDB); err != nil {
			log.Fatalf("ensureSchema: %v", err)
		}
	}

	os.Exit(m.Run())
}

func setDefaultEnv(key, val string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, val)
	}
}

func ensureSchema(db *sql.DB) error {
	schema, err := os.ReadFile("../../../db/init/001_create_users.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(schema))
	return err
}

func setupRepo(t *testing.T) (*database.UserRepository, context.Context) {
	t.Helper()
	if testDB == nil {
		t.Skip("DB not available; start docker compose to run this test")
	}
	ctx := context.Background()
	if _, err := testDB.ExecContext(ctx, "TRUNCATE TABLE users"); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	return database.NewUserRepository(testDB), ctx
}

func TestUserRepository_Save(t *testing.T) {
	repo, ctx := setupRepo(t)

	u, err := user.NewUser("save@example.com", "password123")
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if err := repo.Save(ctx, u); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if u.ID == 0 {
		t.Error("ID should be set after Save")
	}
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	t.Run("存在しないメールアドレスはfalse", func(t *testing.T) {
		repo, ctx := setupRepo(t)

		exists, err := repo.ExistsByEmail(ctx, user.Email("none@example.com"))
		if err != nil {
			t.Fatalf("ExistsByEmail: %v", err)
		}
		if exists {
			t.Error("expected false")
		}
	})

	t.Run("保存済みのメールアドレスはtrue", func(t *testing.T) {
		repo, ctx := setupRepo(t)

		u, _ := user.NewUser("exists@example.com", "password123")
		if err := repo.Save(ctx, u); err != nil {
			t.Fatalf("Save: %v", err)
		}
		exists, err := repo.ExistsByEmail(ctx, user.Email("exists@example.com"))
		if err != nil {
			t.Fatalf("ExistsByEmail: %v", err)
		}
		if !exists {
			t.Error("expected true")
		}
	})
}
