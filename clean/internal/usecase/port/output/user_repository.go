package output

import (
	"context"

	"github.com/kudoutasuku/go-architecture-sample/clean/internal/entity/user"
)

// UserRepository は Interactor が依存する永続化境界。
// Gateway 層で実装される（DIP）。
type UserRepository interface {
	Save(ctx context.Context, u *user.User) error
	ExistsByEmail(ctx context.Context, email user.Email) (bool, error)
}
