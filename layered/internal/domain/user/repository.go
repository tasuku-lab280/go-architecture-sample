package user

import "context"

type Repository interface {
	Save(ctx context.Context, u *User) error
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
