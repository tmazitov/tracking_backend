package order

import (
	"errors"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) CreateOrder(order bl.CreateOrder, isManager bool) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}
	var (
		execString string
		values     []interface{}
	)
	defer s.repo.Close()

	if isManager {
		execString, values = getManagerExecParts(order)
	} else {
		execString, values = getUserExecParts(order)
	}

	_, err = conn.Exec(execString, values...)
	if err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	return nil
}

func getUserExecParts(order bl.CreateOrder) (string, []interface{}) {
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
	return execString, []interface{}{
		order.StartAt,
		order.OwnerID,
		pq.Array(order.PointsID),
		order.Helpers,
		order.Comment,
		order.IsFragileCargo,
	}
}

func getManagerExecParts(order bl.CreateOrder) (string, []interface{}) {
	execString := `insert into orders 
	(
		start_at,
		owner_id,
		manager_id,
		points, 
		helpers, 
		comment_message, 
		is_fragile_cargo
	) 
	values ($1, $2, $3, $4, $5, $6, $7)`
	return execString, []interface{}{
		order.StartAt,
		order.OwnerID,
		order.OwnerID,
		pq.Array(order.PointsID),
		order.Helpers,
		order.Comment,
		order.IsFragileCargo,
	}
}
