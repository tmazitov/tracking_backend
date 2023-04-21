package order

import (
	"errors"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderUpdate(orderId int64, info bl.DB_EditableOrder) error {
	var (
		err error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	defer s.repo.Close()

	execString := `
	UPDATE orders SET 
		start_at=$2,
		helpers=$3,
		points_id=$4,
		comment_message=$5,
		is_fragile_cargo=$6,
		edited_at=now()
	WHERE id=$1 `

	if err = conn.QueryRow(execString, orderId, info.StartAt, info.Helpers, pq.Int64Array(info.PointsID), info.Comment, info.IsFragileCargo).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
