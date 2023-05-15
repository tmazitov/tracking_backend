package order

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/lib/pq"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderList(userId int64, roleId int, filters bl.R_OrderListFilters) ([]bl.DB_Order, error) {

	userFields := []string{
		"owner_id",
		"worker_id",
		"manager_id",
		"",
	}

	orderListItems := []bl.DB_Order{}
	rows, err := s.orderList(userId, userFields[roleId], filters)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	var tempOrderId int64 = 0
	for rows.Next() {
		var (
			ord     bl.DB_Order
			point   bl.Point
			owner   bl.DB_GetUser
			worker  bl.DB_GetUser
			manager bl.DB_GetUser
		)
		err := rows.Scan(
			&ord.ID,
			&ord.Title,
			&ord.CreatedAt,
			&ord.StartAt,
			&ord.EndAt,
			&ord.OrderType,
			&ord.StatusID,
			&ord.Helpers,
			&ord.Comment,
			&ord.IsFragileCargo,
			&ord.IsRegularCustomer,

			&owner.ID,
			&owner.ShortName,
			&owner.RoleID,

			&worker.ID,
			&worker.ShortName,
			&worker.RoleID,

			&manager.ID,
			&manager.ShortName,
			&manager.RoleID,

			&point.ID,
			&point.Floor,
			&point.Title,
			&point.Longitude,
			&point.Latitude,
		)
		if err != nil {
			return nil, errors.New("DB read error: " + err.Error())
		}

		if ord.ID != tempOrderId {
			ord.Owner = owner
			ord.Worker = worker
			ord.Manager = manager
			ord.Points = []bl.Point{}
			ord.Points = append(ord.Points, point)
			orderListItems = append(orderListItems, ord)
			tempOrderId = ord.ID
		} else {
			var listLength int = len(orderListItems) - 1
			orderListItems[listLength].Points = append(orderListItems[listLength].Points, point)
		}
	}
	return orderListItems, err
}

func orderListFiltersToString(filters bl.R_OrderListFilters) (string, []interface{}) {
	var (
		filterString  string = ""
		filterItems   []interface{}
		filterCounter int = 1
	)

	if filters.WorkerId != 0 {
		filterString += fmt.Sprintf("AND worker_id = $%d ", filterCounter)
		filterItems = append(filterItems, filters.WorkerId)
		filterCounter += 1
	}

	if len(filters.Statuses) != 0 {
		fmt.Println(filters.Statuses)
		filterString += fmt.Sprintf("AND status_id = ANY($%d) ", filterCounter)
		filterItems = append(filterItems, pq.Array(filters.Statuses))
		filterCounter += 1
	}

	if len(filters.Types) != 0 {
		filterString += fmt.Sprintf("AND type_id = ANY($%d) ", filterCounter)
		filterItems = append(filterItems, pq.Array(filters.Types))
		filterCounter += 1
	}

	if filters.IsRegularCustomer {
		filterString += fmt.Sprintf("AND is_regular_customer = $%d ", filterCounter)
		filterItems = append(filterItems, filters.IsRegularCustomer)
		filterCounter += 1
	}

	return filterString, filterItems
}

func (s *Storage) orderList(userId int64, roleFieldName string, filters bl.R_OrderListFilters) (*sql.Rows, error) {
	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	var (
		filterString string
		filterItems  []interface{}
	)

	if roleFieldName == "worker_id" {
		filters.WorkerId = userId
	}

	filterString, filterItems = orderListFiltersToString(filters)
	if roleFieldName != "" && roleFieldName != "worker_id" {
		filterString += fmt.Sprintf("AND %s=$%d", roleFieldName, len(filterItems)+1)
		filterItems = append(filterItems, userId)
	}

	execString := `SELECT      		
		orders.id,
		orders.title,
		orders.created_at,
		orders.start_at,       
		orders.end_at,
		orders.type_id,
		orders.status_id,
		orders.helpers,
		orders.comment_message,
		orders.is_fragile_cargo,
		orders.is_regular_customer,

		owner.id,
		owner.short_name,
		owner.role,

		worker.id,
		worker.short_name,
		worker.role,

		manager.id,
		manager.short_name,
		manager.role,

		points.id, 
		points.floor, 
		points.title, 
		ST_X(ST_AsText(points.point)),
		ST_Y(ST_AsText(points.point))
	FROM (
		SELECT * FROM orders WHERE start_at::date=%s::date ` +
		filterString +
		` LIMIT %s OFFSET %s
	) orders
	INNER JOIN points ON points.id=ANY(orders.points_id)
	INNER JOIN users owner ON owner.id=orders.owner_id
	LEFT JOIN users worker ON worker.id=orders.worker_id
	LEFT JOIN users manager ON manager.id=orders.manager_id
	`

	var rowCount uint = bl.DB_OrderListRowCount
	var filtersLen int = len(filterItems)
	execString = fmt.Sprintf(execString,
		"$"+strconv.Itoa(filtersLen+1),
		"$"+strconv.Itoa(filtersLen+2),
		"$"+strconv.Itoa(filtersLen+3),
	)
	filterItems = append(filterItems, filters.Date.Format("2006-01-02"), rowCount, rowCount*filters.Page)

	fmt.Println(execString)

	return conn.Query(execString, filterItems...)
}
