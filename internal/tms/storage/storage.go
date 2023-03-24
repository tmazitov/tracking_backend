package storage

import (
	"log"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	order "github.com/tmazitov/tracking_backend.git/internal/tms/storage/order"
	repo "github.com/tmazitov/tracking_backend.git/pkg/repo"
)

type Storage struct {
	Order bl.OrderStorage
}

func (s *Storage) OrderStorage() bl.OrderStorage {
	return s.Order
}

func NewStorage(repo *repo.Repo, gis *repo.Repo) *Storage {

	storage := Storage{
		Order: order.NewStorage(repo, gis),
	}

	log.Println("tms : storage success")

	return &storage
}
