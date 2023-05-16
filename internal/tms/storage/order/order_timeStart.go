package order

import (
	"errors"
	"time"
)

func (s *Storage) OrderTimeStart(orderId int64) (*time.Time, error) {

	var (
		orderStartFact time.Time
		err            error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE orders 
	SET 
		start_at_fact=TIMEZONE('utc',NOW()),
		status_id=5
	WHERE id=$1
	RETURNING TIMEZONE('utc',NOW())`

	if err = conn.QueryRow(execString, orderId).Scan(&orderStartFact); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return &orderStartFact, nil
}
