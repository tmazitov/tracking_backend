package order

import (
	"errors"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) InsertOrder(order bl.CreateOrder) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `insert into orders 
	(
		start_at, 
		points, 
		helpers, 
		comment_message, 
		is_fragile_cargo
	) 
	values ($1, $2, $3, $4, $5)`
	_, err = conn.Exec(execString, order.StartAt, pq.Array(order.PointsID), order.Helpers, order.Comment, order.IsFragileCargo)
	if err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
