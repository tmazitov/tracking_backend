package order

import (
	"errors"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/app/bl"
)

func (s *Storage) InsertOrder(order bl.CreateOrder) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `insert into orders 
	(
		startAt, 
		endAt, 
		points, 
		helpers, 
		comment_message, 
		is_fragile_cargo
	) 
	values ($1, $2, $3, $4, $5, $6)`
	_, err = conn.Exec(execString, order.StartAt, order.EndAt, pq.Array(order.PointsID), order.Helpers, order.Comment, order.IsFragileCargo)
	if err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
