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
		owner_id,
		points, 
		helpers, 
		comment_message, 
		is_fragile_cargo
	) 
	values ($1, $2, $3, $4, $5, $6)`
	values := []interface{}{
		order.StartAt,
		order.OwnerID,
		pq.Array(order.PointsID),
		order.Helpers,
		order.Comment,
		order.IsFragileCargo,
	}
	_, err = conn.Exec(execString, values...)
	if err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}
