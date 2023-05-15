package bl

import "time"

type OrderStorage interface {
	OrderList(userId int64, roleId int, filters R_OrderListFilters) ([]DB_Order, error)

	OrderGetManagerID(orderId int64) (int64, error)
	OrderGetPointsID(orderId int64) ([]int64, error)
	OrderGetOwnerID(orderId int64) (int64, error)
	OrderGetStatus(orderId int64) (OrderStatus, error)
	OrderGet(orderId int64) (*DB_Order, error)

	OrderUpdateWorker(orderId int64, workerId int64) (*DB_GetUser, error)
	OrderUpdate(orderId int64, info DB_EditableOrder) error

	OrderTimeStart(orderId int64) (*time.Time, error)
	OrderTimeEnd(orderId int64) (*time.Time, error)

	OrderStatusUpgrade(orderId int64) (int, error)
	CreateOrder(order CreateOrder, role UserRole) (int64, error)
	PointsCreate(orderID int64, points []Point) ([]int64, error)
	PointsUpdate(points []Point) ([]int64, error)
	PointsDelete(pointsID []int64) error
	PointsGet(pointsID []int64) ([]Point, error)
}

type UserStorage interface {
	UserGetStaffList() ([]DB_GetUser, error)
	GetUserInfo(userId int64) (DB_GetUser, error)
	UpdateUserShortName(userId int64, shortName string) error
}

type Storage interface {
	OrderStorage() OrderStorage
	UserStorage() UserStorage
}
