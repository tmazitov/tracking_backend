package conductor

import (
	"context"
	"errors"
)

func (c *Conductor) ValidateCode(ctx context.Context, token string) (CodePayload, error) {

	var result CodePayload

	claims, err := c.jwt.ValidateCheck(ctx, token)
	if err != nil {
		return result, errors.New("conductor error: " + err.Error())
	}
	result = CodePayload{
		Email: claims.Email,
		Ip:    claims.Ip,
	}

	return result, nil
}
