package jwt

import (
	"github.com/redis/go-redis/v9"
)

type JwtConfig struct {
	Secret []byte
	Salt   string
}

type JwtStorage struct {
	redis  *redis.Client
	config JwtConfig
}

func NewJwtStorage(config JwtConfig, redis *redis.Client) *JwtStorage {
	return &JwtStorage{
		redis:  redis,
		config: config,
	}
}
