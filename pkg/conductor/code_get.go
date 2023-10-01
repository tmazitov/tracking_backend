package conductor

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

func (c *Conductor) GetCode(ctx context.Context, token string) (string, error) {

	// Create record to check ticket later
	code, err := c.redis.Get(ctx, "che:"+token).Result()
	// no attempts earlier
	if err == redis.Nil {
		return code, errors.New("conductor error: " + ErrInvalidToken.Error())
	}
	if err != nil {
		return code, errors.New("conductor error: " + err.Error())
	}

	return code, nil
}
