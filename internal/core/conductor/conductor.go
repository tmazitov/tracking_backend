package conductor

import (
	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
)

type Config struct {
	Email string
	Pass  string
}

type Conductor struct {
	config Config
	redis  *redis.Client
	jwt    *jwt.JwtStorage

	emailChan chan emailMessage
}

func NewConductor(config Config, jwt *jwt.JwtStorage, redis *redis.Client) Conductor {

	// Ticket chanel
	emailChan := make(chan emailMessage)

	c := Conductor{
		config:    config,
		redis:     redis,
		jwt:       jwt,
		emailChan: emailChan,
	}

	// Setup the worker which send tickets to emails
	go func() {
		c.senderWorker(emailChan)
	}()

	return c
}
