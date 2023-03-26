package bl

type OrderStorage interface {
	OrderList(userId int, roleId int) ([]DB_OrderListItem, error)
	InsertOrder(order CreateOrder) error
	InsertPoint(points []Point) ([]int64, error)
	PointGet(pointsID []int64) ([]Point, error)
}

type UserStorage interface {
	GetUserInfo(userId int) (SelectUserById, error)
}

type Storage interface {
	OrderStorage() OrderStorage
	UserStorage() UserStorage
}
