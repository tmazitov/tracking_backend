package bl

type UserStorage interface {
	CheckUserByEmail(email string) (bool, error)
	CreateUser(email string) (int, error)
}

type Storage interface {
	UserStorage() UserStorage
}
