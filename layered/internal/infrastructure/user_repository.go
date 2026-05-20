package infrastructure

import (
	"context"
	"database/sql"

	"github.com/kudoutasuku/go-architecture-sample/layered/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, u *domain.User) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash) VALUES (?, ?)`,
		u.Email.String(), u.PasswordHash,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email domain.Email) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)`,
		email.String(),
	).Scan(&exists)
	return exists, err
}
