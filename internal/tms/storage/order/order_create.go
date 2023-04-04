package order

import (
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) CreateOrder(order bl.CreateOrder, isManager bool) error {
	conn, err := s.repo.Conn()
	if err != nil {
		return errors.New("DB conn error: " + err.Error())
	}

	tx, err := conn.Begin()
	if err != nil {
		return errors.New("DB transaction begin err: " + err.Error())
	}

	var (
		orderID            int64
		execString         string
		orderValues        []interface{}
		orderToPointValues []interface{}
	)
	defer s.repo.Close()

	if isManager {
		execString, orderValues = getManagerExecParts(order)
	} else {
		execString, orderValues = getUserExecParts(order)
	}

	if err = tx.QueryRow(execString, orderValues...).Scan(&orderID); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	var counter int = 1
	var unknownValuesString string = ""
	for index, pointID := range order.PointsID {
		orderToPointValues = append(orderToPointValues, orderID, pointID)
		unknownValuesString += fmt.Sprintf("($%d, $%d)", counter, counter+1)
		if index != len(order.PointsID)-1 {
			unknownValuesString += ", "
		}
		counter += 2
	}

	execString = fmt.Sprintf(`INSERT INTO points_to_orders (order_id, point_id) VALUES %s`, unknownValuesString)
	if err = tx.QueryRow(execString, orderToPointValues...).Err(); err != nil {
		return errors.New("DB exec error: " + err.Error())
	}

	if err = tx.Commit(); err != nil {
		return errors.New("DB transaction continue err: " + err.Error())
	}

	return nil
}

func getUserExecParts(order bl.CreateOrder) (string, []interface{}) {
	execString := `INSERT INTO orders 
	(
		start_at,
		owner_id,
		points_id, 
		helpers, 
		comment_message, 
		is_fragile_cargo
	) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`
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
	execString := `INSERT INTO orders 
	(
		start_at,
		owner_id,
		manager_id,
		points_id, 
		helpers, 
		comment_message, 
		is_fragile_cargo
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`
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
