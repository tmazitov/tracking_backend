package ws

import (
	"context"
	"log"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
)

func (c *Client) checkAccess(jwt *jwt.JwtStorage) (*jwt.AccessClaims, error) {

	userPayload, err := c.authMiddleware(jwt, c.access)
	if !c.isAuthorized && c.isRefreshNeeded {
		jsonReport, err := NewResponse(401, "Unauthorized").Marshal()
		if err != nil {
			log.Println("wrong report marshal")
			return nil, err
		}

		c.send <- jsonReport
	}

	return userPayload, err
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
	c.role = bl.UserRole(userPayload.RoleId)
	c.access = access

	return userPayload, nil
}
