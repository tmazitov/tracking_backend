package storage

import (
	"log"

	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	"github.com/tmazitov/tracking_backend.git/internal/admin/storage/user"
	repo "github.com/tmazitov/tracking_backend.git/pkg/repo"
)

type Storage struct {
	User bl.UserStorage
}

func (s *Storage) UserStorage() bl.UserStorage {
	return s.User
}

func NewStorage(repo *repo.Repo) *Storage {

	storage := Storage{
		User: user.NewStorage(repo),
	}

	log.Println("admin : storage success")

	return &storage
}
