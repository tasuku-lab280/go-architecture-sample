package input

import "context"

type RegisterUserInputData struct {
	Email    string
	Password string
}

// RegisterUserInputPort は Controller が依存するユースケースの入力境界。
type RegisterUserInputPort interface {
	Handle(ctx context.Context, in RegisterUserInputData) error
}
