package conductor

import "context"

func (c *Conductor) DeleteCode(ctx context.Context, codeToken string) error {
	return c.redis.Del(ctx, "che:"+codeToken).Err()
}
