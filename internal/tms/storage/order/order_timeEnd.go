package order

import (
	"errors"
	"time"
)

func (s *Storage) OrderTimeEnd(orderId int64) (*time.Time, error) {

	var (
		orderEndFact time.Time
		err          error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE orders 
	SET 
		end_at_fact=NOW(),
		status_id=1
	WHERE id=$1
	RETURNING NOW()`

	if err = conn.QueryRow(execString, orderId).Scan(&orderEndFact); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return &orderEndFact, nil
}
