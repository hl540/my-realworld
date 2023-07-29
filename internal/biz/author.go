package biz

type Author struct {
	ID        uint
	Username  string
	Image     string
	Bio       string
	Following bool
}

type AuthorUseCase struct {
}
