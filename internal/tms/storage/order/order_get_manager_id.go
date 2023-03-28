package order

import (
	"database/sql"
	"errors"
)

func (s *Storage) OrderGetManagerID(orderId int) (sql.NullInt64, error) {

	var (
		execString string
		result     sql.NullInt64
		err        error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return result, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString = `SELECT manager_id FROM orders WHERE id=$1`

	if err = conn.QueryRow(execString, orderId).Scan(&result); err != nil {
		return result, errors.New("DB exec error: " + err.Error())
	}
	return result, err
}
