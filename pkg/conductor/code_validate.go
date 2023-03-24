package conductor

import "context"

func (c *Conductor) ValidateCode(ctx context.Context, token string) (CodePayload, error) {

	var result CodePayload

	claims, err := c.jwt.ValidateCheck(ctx, token)
	if err != nil {
		return result, err
	}
	result = CodePayload{
		Email: claims.Email,
		Ip:    claims.Ip,
	}

	return result, nil
}
