package user

import (
	repo "github.com/tmazitov/tracking_backend.git/internal/core/repo"
)

type Storage struct {
	repo *repo.Repo
}

func NewStorage(repo *repo.Repo) *Storage {
	return &Storage{repo: repo}
}
