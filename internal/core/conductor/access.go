package conductor

import (
	"context"

	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
)

func (c *Conductor) CreateTokenPair(ctx context.Context, claims jwt.AccessClaims) (jwt.TokenPair, error) {
	return c.jwt.CreateTokenPair(ctx, claims)
}

func (c *Conductor) DeleteTokenPair(ctx context.Context, access string, refresh string) error {
	return c.jwt.DeleteTokenPair(ctx, access, refresh)
}

func (c *Conductor) ValidateRefresh(ctx context.Context, refresh string) (*jwt.AccessClaims, error) {
	return c.jwt.ValidateRefresh(ctx, refresh)
}
