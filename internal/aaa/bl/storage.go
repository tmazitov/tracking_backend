package bl

type UserStorage interface {
	CheckUserByEmail(email string) (bool, error)
	CreateUser(email string) error
}

type Storage interface {
	UserStorage() UserStorage
}
