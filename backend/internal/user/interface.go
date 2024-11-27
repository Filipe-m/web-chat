package user

type RepositoryInterface interface {
	GetUserByID(user *User) error
	Create(user *User) error
	GetUserByEmail(email string) (User, error)
}
