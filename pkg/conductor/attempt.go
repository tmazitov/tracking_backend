package conductor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *Conductor) checkAuthAttempts(ctx context.Context, email string) (int, error) {
	// Create record to check ticket later
	val, err := c.redis.Get(ctx, "auth_attempt:"+email).Result()

	// no attempts earlier
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, errors.New("conductor error: " + err.Error())
	}

	attemptCount, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("conductor error: " + err.Error())
	}

	// too many attempts
	if attemptCount >= MAX_ATTEND_COUNT {
		return attemptCount, ErrTooManyAttempts
	}

	return attemptCount, nil
}

func (c *Conductor) updateAuthAttempts(ctx context.Context, email string, count int) error {
	return c.redis.Set(ctx, "auth_attempt:"+email, count, 5*time.Minute).Err()
}
