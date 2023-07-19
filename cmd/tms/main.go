package main

import (
	"flag"

	"github.com/redis/go-redis/v9"
	config "github.com/tmazitov/tracking_backend.git/config/tms"
	rest "github.com/tmazitov/tracking_backend.git/internal/tms/rest"
	storage "github.com/tmazitov/tracking_backend.git/internal/tms/storage"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	"github.com/tmazitov/tracking_backend.git/pkg/repo"
)

func configPath() *string {
	path := flag.String("config", "", "Path to the config.json file")
	flag.Parse()
	return path
}

func main() {

	config := config.Config{Path: *configPath()}
	if err := config.Setup(); err != nil {
		panic(err)
	}

	storeConf := config.RepoConfig()
	store := &repo.Repo{Config: storeConf}

	gisConf := config.GisConfig()
	gis := &repo.Repo{Config: gisConf}

	// Setup redis api
	redis := redis.NewClient(config.RedisConfig())

	// Setup jwt api
	jwtConf := config.JwtConfig()
	jwt := jwt.NewJwtStorage(jwtConf, redis)

	storage := storage.NewStorage(store, store, gis)

	router := rest.NewRouter("tms", redis, storage, jwt)
	router.Run("5001")
}
