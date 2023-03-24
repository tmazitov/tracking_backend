package bl

type OrderStorage interface {
	InsertOrder(order CreateOrder) error
	InsertPoint(points []Point) ([]int64, error)
}

type UserStorage interface {
	GetUSer(int)
}

type Storage interface {
	OrderStorage() OrderStorage
}
