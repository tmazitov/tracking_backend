package main

import (
	config "github.com/tmazitov/tracking_backend.git/config/aaa"
	rest "github.com/tmazitov/tracking_backend.git/internal/aaa/rest"
	storage "github.com/tmazitov/tracking_backend.git/internal/aaa/storage"
	"github.com/tmazitov/tracking_backend.git/internal/core/repo"
)

func main() {
	config := config.Config{Path: "../../config/aaa/config.json"}
	if err := config.Setup(); err != nil {
		panic(err)
	}

	storeConf := config.RepoConfig()
	store := &repo.Repo{Config: storeConf}

	storage := storage.NewStorage(store)

	router := rest.NewRouter("/aaa/api", storage)
	router.Run("5000")
}
