package main

import (
	"github.com/redis/go-redis/v9"
	config "github.com/tmazitov/tracking_backend.git/config/aaa"
	rest "github.com/tmazitov/tracking_backend.git/internal/aaa/rest"
	storage "github.com/tmazitov/tracking_backend.git/internal/aaa/storage"
	"github.com/tmazitov/tracking_backend.git/internal/core/conductor"
	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
	"github.com/tmazitov/tracking_backend.git/internal/core/repo"
)

func main() {
	config := config.Config{Path: "../../config/aaa/config.json"}
	if err := config.Setup(); err != nil {
		panic(err)
	}

	storeConf := config.RepoConfig()
	store := &repo.Repo{Config: storeConf}

	redis := redis.NewClient(config.RedisConfig())

	jwtConf := config.JwtConfig()
	jwt := jwt.NewJwtStorage(jwtConf, redis)

	condConf := config.CondConfig()
	conductor := conductor.NewConductor(condConf, jwt, redis)

	storage := storage.NewStorage(store)

	router := rest.NewRouter("/aaa/api", storage, conductor)
	router.Run("5000")
}
