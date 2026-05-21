package output

type RegisterUserOutputData struct {
	ID    int64
	Email string
}

// RegisterUserPresenter は Interactor が依存する出力境界。
// Presenter 層で実装され、ViewModel への変換を担う。
type RegisterUserPresenter interface {
	Present(out RegisterUserOutputData)
	PresentError(err error)
}
