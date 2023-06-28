package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/tmazitov/tracking_backend.git/internal/admin/rest"
	storage "github.com/tmazitov/tracking_backend.git/internal/admin/storage"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	"github.com/tmazitov/tracking_backend.git/pkg/repo"

	config "github.com/tmazitov/tracking_backend.git/config/admin"
)

func main() {

	// Setup service config
	config := config.Config{Path: "../../config/admin/config.json"}
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

	// Setup GRPC
	// grpc.SetupServer(":5100")

	// Setup router
	router := rest.NewRouter("adm", storage, jwt)
	router.Run("5002")
}
