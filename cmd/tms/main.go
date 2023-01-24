package main

import (
	config "github.com/tmazitov/tracking_backend.git/config/tms"
	"github.com/tmazitov/tracking_backend.git/internal/core/repo"
	rest "github.com/tmazitov/tracking_backend.git/internal/tms/rest"
	storage "github.com/tmazitov/tracking_backend.git/internal/tms/storage"
)

func main() {

	config := config.Config{Path: "../../config/tms/config.json"}
	if err := config.Setup(); err != nil {
		panic(err)
	}

	storeConf := config.RepoConfig()
	store := &repo.Repo{Config: storeConf}

	gisConf := config.GisConfig()
	gis := &repo.Repo{Config: gisConf}

	storage := storage.NewStorage(store, gis)

	router := rest.NewRouter("/tms/api", storage)
	router.Run("5001")
}
