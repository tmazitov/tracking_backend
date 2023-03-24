package conductor

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (c *Conductor) GetTicket(ctx context.Context, token string) (Ticket, error) {

	// Create record to check ticket later
	tiketData, err := c.redis.Get(ctx, "che:"+token).Result()
	ticket := Ticket{Code: tiketData}

	// no attempts earlier
	if err == redis.Nil {
		return ticket, ErrInvalidToken
	}

	if err != nil {
		return ticket, err
	}

	return ticket, nil
}
