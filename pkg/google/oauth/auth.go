package oauth

import (
	"context"
	"google.golang.org/api/idtoken"
)

type GoogleAuth struct {
	ClientID string
}

func NewGoogleAuth(clientId string) *GoogleAuth {
	return &GoogleAuth{ClientID: clientId}
}

func (g *GoogleAuth) Verify(ctx context.Context, idToken string) (*idtoken.Payload, error) {
	return idtoken.Validate(ctx, idToken, g.ClientID)
}
