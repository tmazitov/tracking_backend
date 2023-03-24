package order

import (
	repo "github.com/tmazitov/tracking_backend.git/pkg/repo"
)

type Storage struct {
	repo *repo.Repo
	gis  *repo.Repo
}

func NewStorage(repo *repo.Repo, gis *repo.Repo) *Storage {
	return &Storage{repo: repo, gis: gis}
}
