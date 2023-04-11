package token

import (
	"context"
	"fmt"

	"github.com/osalomon89/test-crud-api/internal/platform/jwt"
)

type jwtToken struct {
	jwt *jwt.JWT
}

func NewTokenGenerator(jwt *jwt.JWT) *jwtToken {
	return &jwtToken{jwt: jwt}
}

func (j *jwtToken) Get(ctx context.Context, id uint, email string) (string, error) {
	return j.jwt.GenerateToken(id, email)
}

func (j *jwtToken) Validate(ctx context.Context, token string) (*Info, error) {
	accesTokenInfo, err := j.jwt.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("token validation error: %w", err)
	}

	return &Info{
		ID:    accesTokenInfo.ID,
		Email: accesTokenInfo.Email,
	}, nil
}
