package order

import (
	"database/sql"
	"errors"
)

func (s *Storage) OrderGetManagerID(orderId int) (int64, error) {

	var (
		execString     string
		orderManagerID sql.NullInt64
		err            error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return -1, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `SELECT manager_id FROM orders WHERE id=$1`

	if err = conn.QueryRow(execString, orderId).Scan(&orderManagerID); err != nil {
		return -1, errors.New("DB exec error: " + err.Error())
	}
	return orderManagerID.Int64, err
}
