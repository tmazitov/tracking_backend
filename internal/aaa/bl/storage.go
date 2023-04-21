package bl

type UserStorage interface {
	CheckUserByEmail(email string) (int64, int, error)
	CreateUser(email string) (int64, error)
}

type Storage interface {
	UserStorage() UserStorage
}
