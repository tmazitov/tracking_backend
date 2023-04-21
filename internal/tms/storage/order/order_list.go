package order

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderList(userId int64, roleId int, filters bl.R_OrderListFilters) ([]bl.DB_OrderListItem, error) {

	userFields := []string{
		"owner_id",
		"worker_id",
		"manager_id",
		"",
	}

	result := []bl.DB_OrderListItem{}
	rows, err := s.orderList(userId, userFields[roleId], filters)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}
	for rows.Next() {
		ord := bl.DB_OrderListItem{}
		err := rows.Scan(
			&ord.ID,
			&ord.Title,
			&ord.CreatedAt,
			&ord.StartAt,
			&ord.EndAt,
			&ord.OwnerID,
			&ord.WorkerID,
			&ord.ManagerID,
			&ord.OrderType,
			&ord.StatusID,
			&ord.PointsID,
			&ord.Helpers,
			&ord.Comment,
			&ord.IsFragileCargo,
			&ord.IsRegularCustomer,
		)
		if err != nil {
			return nil, errors.New("DB read error: " + err.Error())
		}
		result = append(result, ord)
	}
	return result, err
}

func orderListFiltersToString(filters bl.R_OrderListFilters) (string, []interface{}) {
	var (
		filterString  string = ""
		filterItems   []interface{}
		filterCounter int = 1
	)

	if filters.WorkerId != 0 {
		filterString += fmt.Sprintf("AND worker_id = $%d", filterCounter)
		filterItems = append(filterItems, filters.WorkerId)
		filterCounter += 1
	}

	if filters.StatusID != 0 {
		filterString += fmt.Sprintf("AND status_id = $%d", filterCounter)
		filterItems = append(filterItems, filters.StatusID)
		filterCounter += 1
	}

	if filters.TypeId != 0 {
		filterString += fmt.Sprintf("AND type_id = $%d", filterCounter)
		filterItems = append(filterItems, filters.TypeId)
		filterCounter += 1
	}

	if filters.IsRegularCustomer {
		filterString += fmt.Sprintf("AND is_regular_customer = $%d", filterCounter)
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
		id,
		title,
		created_at,
		start_at,       
		end_at,
		owner_id,
		worker_id,   
		manager_id,
		type_id,
		status_id,
		points_id,
		helpers,
		comment_message,
		is_fragile_cargo,
		is_regular_customer
	FROM orders
	WHERE start_at::date='%s'::date `

	var rowCount uint = bl.DB_OrderListRowCount
	execString = execString + filterString
	execString = fmt.Sprintf(execString+"	LIMIT %d OFFSET %d", filters.Date.Format("2006-01-02"), rowCount, rowCount*filters.Page)

	fmt.Println(execString)

	return conn.Query(execString, filterItems...)
}
