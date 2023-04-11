package token

import "context"

type Info struct {
	ID    uint
	Email string
}

type Token struct {
	Access string
	Info   Info
}

type Service interface {
	Get(ctx context.Context, id uint, email string) (string, error)
	Validate(ctx context.Context, token string) (*Info, error)
}
