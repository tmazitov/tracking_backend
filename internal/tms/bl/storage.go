package bl

type OrderStorage interface {
	OrderList(userId int, roleId int) ([]DB_OrderListItem, error)
	OrderGetManagerID(orderId int) (int64, error)
	OrderGetPointsID(orderId int) ([]int64, error)
	OrderGetOwnerID(orderId int) (int64, error)
	OrderGetStatus(orderId int) (OrderStatus, error)
	OrderUpdate(orderId int, info DB_EditableOrder) error
	OrderStatusUpgrade(orderId int) (int, error)
	CreateOrder(order CreateOrder, isManager bool) (int64, error)
	PointsCreate(orderID int, points []Point) ([]int64, error)
	PointsUpdate(points []Point) ([]int64, error)
	PointsDelete(pointsID []int64) error
	PointsGet(pointsID []int64) ([]Point, error)
}

type UserStorage interface {
	UserGetStaffList() ([]DB_GetUser, error)
	GetUserInfo(userId int) (DB_GetUser, error)
	UpdateUserShortName(userId int, shortName string) error
}

type Storage interface {
	OrderStorage() OrderStorage
	UserStorage() UserStorage
}
