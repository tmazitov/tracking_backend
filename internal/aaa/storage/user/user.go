package user

import (
	repo "github.com/tmazitov/tracking_backend.git/pkg/repo"
)

type Storage struct {
	repo *repo.Repo
}

func NewStorage(repo *repo.Repo) *Storage {
	return &Storage{repo: repo}
}
