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
		end_at=$3,
		title=$4,
		comment_message=$5,
		points_id=$6,
		type_id=$7,
		is_regular_customer=$8,
		edited_at=now()
	WHERE id=$1 `

	if err = conn.QueryRow(execString,
		orderId,
		info.StartAt,
		info.EndAt,
		info.Title,
		info.Comment,
		pq.Int64Array(info.PointsID),
		info.OrderType,
		info.IsRegularCustomer,
	).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
