package bl

type OrderStorage interface {
	InsertOrder(order CreateOrder) error
	InsertPoint(points []Point) ([]int64, error)
}

type UserStorage interface {
	GetUserInfo(userId int) (SelectUserById, error)
}

type Storage interface {
	OrderStorage() OrderStorage
	UserStorage() UserStorage
}
