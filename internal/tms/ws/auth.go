package ws

import (
	"context"

	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

func (c *Client) checkAccess(jwt *jwt.JwtStorage) (*jwt.AccessClaims, error) {
	return c.authMiddleware(jwt, c.access)
}

func (c *Client) authMiddleware(jwt *jwt.JwtStorage, access string) (*jwt.AccessClaims, error) {
	var ctx context.Context = context.Background()

	userPayload, err := jwt.ValidateAccessToken(ctx, access)
	if err != nil {
		c.isAuthorized = false

		return nil, err
	}

	c.isAuthorized = true

	c.userId = userPayload.UserId
	c.access = access

	return userPayload, nil
}
