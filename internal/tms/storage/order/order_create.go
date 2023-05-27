package order

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) CreateOrder(order bl.CreateOrder, role bl.UserRole) (int64, error) {

	var (
		orderToPointValues  []interface{}
		createOrderValues   []interface{}
		unknownValuesString string
		pointsID            []int64
		orderID             int64
		execString          string
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
	var (
		queryString            string
		queryFieldsSpotsString string
		queryFieldsString      string
		queryFieldsValues      []interface{}
	)
	fmt.Println(order.OrderType)
	queryFieldsString, queryFieldsSpotsString, queryFieldsValues = getQueryItems(role, order, pointsID)
	fmt.Println(queryFieldsString)
	fmt.Println(queryFieldsSpotsString)
	fmt.Println(queryFieldsValues...)
	queryString = fmt.Sprintf(`INSERT INTO orders ( %s ) 
	VALUES ( %s )
	RETURNING id`, queryFieldsString, queryFieldsSpotsString)

	if err = tx.QueryRow(queryString, queryFieldsValues...).Scan(&orderID); err != nil {
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

func getQueryItems(role bl.UserRole, order bl.CreateOrder, pointsID []int64) (string, string, []interface{}) {
	var (
		querySpots  []string
		queryFields []string
		queryItems  []interface{}
	)

	queryFields = []string{
		"start_at",
		"owner_id",
		"points_id",
		"title",
		"helpers",
		"comment_message",
		"type_id",
		"is_fragile_cargo",
	}

	queryItems = []interface{}{
		order.StartAt,
		order.OwnerID,
		pq.Int64Array(pointsID),
		order.Title,
		order.Helpers,
		order.Comment,
		order.OrderType,
		order.IsFragileCargo,
	}

	querySpots = []string{"$1", "$2", "$3", "$4", "$5", "$6", "$7", "$8"}

	if role == bl.Admin || role == bl.Manager {
		queryFields = append(queryFields, "is_regular_customer", "manager_id")
		queryItems = append(queryItems, order.IsRegularCustomer, order.OwnerID)
		querySpots = append(querySpots, "$9", "$10")
		fmt.Println("worker id is ", order.WorkerID)
		if order.WorkerID != 0 {
			queryFields = append(queryFields, "worker_id", "status_id")
			queryItems = append(queryItems, order.WorkerID, 4)
			querySpots = append(querySpots, "$11", "$12")
		}
		if !order.EndAt.IsZero() {
			queryFields = append(queryFields, "end_at")
			queryItems = append(queryItems, order.EndAt)
			querySpots = append(querySpots, fmt.Sprintf("$%d", len(querySpots)+1))
		}
	}

	return strings.Join(queryFields, ", "), strings.Join(querySpots, ", "), queryItems
}
