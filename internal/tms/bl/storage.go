package bl

import "database/sql"

type OrderStorage interface {
	OrderList(userId int, roleId int) ([]DB_OrderListItem, error)
	OrderGetManagerID(orderId int) (sql.NullInt64, error)
	OrderStatusUpgrade(orderId int) (int, error)
	CreateOrder(order CreateOrder, isManager bool) error
	CreatePoints(points []Point) ([]int64, error)
	GetPoints(pointsID []int64) ([]Point, error)
}

type UserStorage interface {
	GetUserInfo(userId int) (SelectUserById, error)
	UpdateUserShortName(userId int, shortName string) error
}

type Storage interface {
	OrderStorage() OrderStorage
	UserStorage() UserStorage
}
