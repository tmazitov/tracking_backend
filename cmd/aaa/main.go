package main

import (
	"github.com/redis/go-redis/v9"
	config "github.com/tmazitov/tracking_backend.git/config/aaa"
	rest "github.com/tmazitov/tracking_backend.git/internal/aaa/rest"
	storage "github.com/tmazitov/tracking_backend.git/internal/aaa/storage"
	"github.com/tmazitov/tracking_backend.git/pkg/conductor"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	"github.com/tmazitov/tracking_backend.git/pkg/repo"
)

func main() {

	// Setup service config
	config := config.Config{Path: "../../config/aaa/config.json"}
	if err := config.Setup(); err != nil {
		panic(err)
	}

	// Setup database storage
	storeConf := config.RepoConfig()
	store := &repo.Repo{Config: storeConf}
	storage := storage.NewStorage(store)

	// Setup redis api
	redis := redis.NewClient(config.RedisConfig())

	// Setup jwt api
	jwtConf := config.JwtConfig()
	jwt := jwt.NewJwtStorage(jwtConf, redis)

	// Setup conductor api (for authorization with codes and tokens)
	condConf := config.CondConfig()
	conductor := conductor.NewConductor(condConf, jwt, redis)

	// Setup GRPC
	// grpc.SetupServer(":5100")

	// Setup router
	router := rest.NewRouter("/aaa/api", storage, conductor, jwt)
	router.Run("5000")
}
