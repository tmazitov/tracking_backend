package bl

type UserStorage interface {
	CheckUserByEmail(email string) (int, int, error)
	CreateUser(email string) (int, error)
}

type Storage interface {
	UserStorage() UserStorage
}
