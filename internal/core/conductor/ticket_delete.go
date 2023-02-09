package conductor

import "context"

func (c *Conductor) DeleteTicket(ctx context.Context, tiketToken string) error {
	return c.redis.Del(ctx, "che:"+tiketToken).Err()
}
