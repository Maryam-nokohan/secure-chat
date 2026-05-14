package user

type UserRepository interface {

	CreateUser(user User) error
	FindUser(userid string) error
	EditUser(user User) error
	DeleteUser(user User) error
}