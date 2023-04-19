package order

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderList(userId int, roleId int) ([]bl.DB_OrderListItem, error) {

	userFields := []string{
		"owner_id",
		"worker_id",
		"manager_id",
		"",
	}

	result := []bl.DB_OrderListItem{}
	rows, err := s.orderList(userId, userFields[roleId])
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

func (s *Storage) orderList(userId int, roleFieldName string) (*sql.Rows, error) {
	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT      		
		id,
		title,
		created_at,
		start_at,       
		end_at,
		owner_id,
		worker_id,   
		manager_id,
		status_id,
		points_id,
		helpers,
		comment_message,
		is_fragile_cargo,
		is_regular_customer
	FROM orders`

	if roleFieldName == "" {
		return conn.Query(execString)
	}

	execString += fmt.Sprintf(" WHERE %s=$1", roleFieldName)
	return conn.Query(execString, userId)
}
