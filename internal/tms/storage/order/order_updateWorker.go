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
		worker_id=$2,
		status_id=4
	FROM (
		SELECT
			*
		FROM orders
		WHERE id=$1
		FOR UPDATE
	) o
		INNER JOIN users w ON w.id=$2
	RETURNING w.id, w.short_name, w.role`

	if err = conn.QueryRow(execString, orderId, workerId).Scan(&worker.ID, &worker.ShortName, &worker.RoleID); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return &worker, nil
}
