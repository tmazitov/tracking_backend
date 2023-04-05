package order

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) CreateOrder(order bl.CreateOrder, isManager bool) (int64, error) {

	var (
		orderToPointValues  []interface{}
		createOrderValues   []interface{}
		unknownValuesString string
		pointsID            []int64
		orderID             int64
		execString          string
		orderValues         []interface{}
		tx                  *sql.Tx
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return 0, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	tx, err = conn.Begin()
	if err != nil {
		return 0, errors.New("DB transaction begin error: " + err.Error())
	}

	// Create initial points of order and get the points id
	createOrderValues, unknownValuesString = getCreatablePointsValues(order.Points)
	execString = fmt.Sprintf(`INSERT INTO points 
		(title, step_id, floor, point) 
		VALUES %s
		RETURNING id`, unknownValuesString)

	rows, err := conn.Query(execString, createOrderValues...)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("DB exec error: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return 0, errors.New("DB get data error: " + err.Error())
		}
		pointsID = append(pointsID, id)
	}

	// Create new order and get the new order id
	if isManager {
		execString, orderValues = getManagerExecParts(order, pointsID)
	} else {
		execString, orderValues = getUserExecParts(order, pointsID)
	}

	if err = tx.QueryRow(execString, orderValues...).Scan(&orderID); err != nil {
		tx.Rollback()
		return 0, errors.New("DB exec error: " + err.Error())
	}

	// Create relationship between the points and order
	var counter int = 1
	unknownValuesString = ""
	for index, pointID := range pointsID {
		orderToPointValues = append(orderToPointValues, orderID, pointID)
		unknownValuesString += fmt.Sprintf("($%d, $%d)", counter, counter+1)
		if index != len(pointsID)-1 {
			unknownValuesString += ", "
		}
		counter += 2
	}

	execString = fmt.Sprintf(`INSERT INTO points_to_orders (order_id, point_id) VALUES %s`, unknownValuesString)
	if err = tx.QueryRow(execString, orderToPointValues...).Err(); err != nil {
		tx.Rollback()
		return 0, errors.New("DB exec error: " + err.Error())
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, errors.New("DB transaction continue err: " + err.Error())
	}

	return orderID, nil
}

func getUserExecParts(order bl.CreateOrder, pointsID []int64) (string, []interface{}) {
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
		pq.Int64Array(pointsID),
		order.Helpers,
		order.Comment,
		order.IsFragileCargo,
	}
}

func getManagerExecParts(order bl.CreateOrder, pointsID []int64) (string, []interface{}) {
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
		pq.Int64Array(pointsID),
		order.Helpers,
		order.Comment,
		order.IsFragileCargo,
	}
}
