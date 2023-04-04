package order

import (
	"errors"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderUpdateMainInfo(orderId int, info bl.R_EditableOrder) (pq.Int64Array, error) {
	var (
		pointsID pq.Int64Array
		err      error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE orders SET 
		start_at=$2,
		helpers=$3,
		comment=$4,
		is_fragile_cargo=$5,
		edited_at=now()
	WHERE id=$1 RETURNING points`

	if err = conn.QueryRow(execString, orderId, info.StartAt, info.Helpers, info.IsFragileCargo).Scan(&pointsID); err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	return pointsID, nil
}
