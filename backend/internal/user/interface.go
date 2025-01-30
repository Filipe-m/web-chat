package user

type RepositoryInterface interface {
	GetUserByID(user *User) error
	Create(user *User) error
	GetUserByEmail(email string) (User, error)
}

type ServiceInterface interface {
    Create(User) error
    Authenticate(User) (string, error)
    GetUserByID(int) (User, error)
}