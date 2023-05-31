package order

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderUpdateWorker(orderId int64, workerId int64) (*bl.DB_GetUser, error) {
	var (
		worker bl.DB_GetUser
		err    error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE orders
	SET 
		worker_id = $2,
		status_id = $3
	FROM (
		SELECT * FROM users WHERE id = $2
	) worker
	WHERE orders.id = $1
	RETURNING worker.id, worker.short_name, worker.role`

	if err = conn.QueryRow(execString, orderId, workerId, bl.OrderAccepted).Scan(&worker.ID, &worker.ShortName, &worker.RoleID); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return &worker, nil
}
