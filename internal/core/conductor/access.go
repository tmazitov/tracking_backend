package conductor

import (
	"context"

	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
)

func (c *Conductor) CreateAccess(ctx context.Context, claims jwt.AccessClaims) (jwt.TokenPair, error) {
	return c.jwt.CreateTokenPair(ctx, claims)
}
