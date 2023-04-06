package order

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderGetStatus(orderId int) (bl.OrderStatus, error) {
	var (
		execString  string
		orderStatus bl.OrderStatus
		err         error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return -1, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `SELECT points_id FROM orders WHERE id=$1`

	if err = conn.QueryRow(execString, orderId).Scan(&orderStatus); err != nil {
		return -1, errors.New("DB exec error: " + err.Error())
	}
	return orderStatus, err
}
