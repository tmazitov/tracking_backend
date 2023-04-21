package order

import (
	"errors"

	"github.com/lib/pq"
)

func (s *Storage) OrderGetPointsID(orderId int64) ([]int64, error) {

	var (
		execString string
		pointsID   pq.Int64Array
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `SELECT points_id FROM orders WHERE id=$1`

	if err = conn.QueryRow(execString, orderId).Scan(&pointsID); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}
	return pointsID, err
}
