package storage

import (
	"log"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	order "github.com/tmazitov/tracking_backend.git/internal/tms/storage/order"
	user "github.com/tmazitov/tracking_backend.git/internal/tms/storage/user"
	repo "github.com/tmazitov/tracking_backend.git/pkg/repo"
)

type Storage struct {
	Order bl.OrderStorage
	User  bl.UserStorage
}

func (s *Storage) OrderStorage() bl.OrderStorage {
	return s.Order
}

func (s *Storage) UserStorage() bl.UserStorage {
	return s.User
}

func NewStorage(userRepo *repo.Repo, orderRepo *repo.Repo, gis *repo.Repo) *Storage {

	storage := Storage{
		Order: order.NewStorage(orderRepo, gis),
		User:  user.NewStorage(userRepo),
	}

	log.Println("tms : storage success")

	return &storage
}
