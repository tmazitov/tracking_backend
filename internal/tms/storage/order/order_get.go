package order

import (
	"errors"

	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

func (s *Storage) OrderGet(orderID int64) (*bl.DB_Order, error) {
	var (
		order   bl.DB_Order
		bill    bl.DB_OrderBill
		point   bl.Point
		owner   bl.DB_GetUser
		worker  bl.DB_GetUser
		manager bl.DB_GetUser
		err     error
	)

	conn, err := s.repo.Conn()
	if err != nil {
		return nil, errors.New("DB conn error: " + err.Error())
	}

	defer s.repo.Close()

	execString := `SELECT      		
		orders.id,
		orders.title,
		orders.created_at,
		orders.start_at,       
		orders.end_at,
		orders.type_id,
		orders.status_id,
		orders.comment_message,
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
		ST_Y(ST_AsText(points.point)),

		bill.car_type_id,
		bill.helper_count,
		bill.helper_hours,
		bill.is_fragile_cargo
	FROM (
		SELECT * FROM orders WHERE id=$1
	) orders
	INNER JOIN points ON points.id=ANY(orders.points_id)
	INNER JOIN users owner ON owner.id=orders.owner_id
	INNER JOIN order_bills bill ON bill.order_id=orders.id
	LEFT JOIN users worker ON worker.id=orders.worker_id
	LEFT JOIN users manager ON manager.id=orders.manager_id
	`

	rows, err := conn.Query(execString, orderID)
	if err != nil {
		return nil, errors.New("DB exec error: " + err.Error())
	}

	var tempOrderId int64 = 0
	for rows.Next() {
		var ()
		err := rows.Scan(
			&order.ID,
			&order.Title,
			&order.CreatedAt,
			&order.StartAt,
			&order.EndAt,
			&order.OrderType,
			&order.StatusID,
			&order.Comment,
			&order.IsRegularCustomer,

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

			&bill.CarTypeID,
			&bill.HelperCount,
			&bill.HelperHours,
			&bill.IsFragileCargo,
		)
		if err != nil {
			return nil, errors.New("DB read error: " + err.Error())
		}

		if order.ID != tempOrderId {
			order.Owner = owner
			order.Worker = worker
			order.Manager = manager
			order.Points = []bl.Point{}
			tempOrderId = order.ID
		}

		order.Points = append(order.Points, point)
	}

	return &order, err
}
