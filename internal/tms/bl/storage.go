package bl

type OrderStorage interface {
	OrderList(userId int, roleId int) ([]DB_OrderListItem, error)
	CreateOrder(order CreateOrder, isManager bool) error
	CreatePoints(points []Point) ([]int64, error)
	GetPoints(pointsID []int64) ([]Point, error)
}

type UserStorage interface {
	GetUserInfo(userId int) (SelectUserById, error)
}

type Storage interface {
	OrderStorage() OrderStorage
	UserStorage() UserStorage
}
