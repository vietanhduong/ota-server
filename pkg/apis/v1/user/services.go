package auth

import (
	"context"
)

type service struct {
}

//var secret = env.GetEnvAsStringOrFallback("SECRET", "some-thing-very-secret")

func NewService() *service {
	return &service{
	}
}

func (s *service) Login(ctx context.Context, idToken string) (*Token, error) {

	return nil, nil
}

func (s *service) GenerateToken() {

}
