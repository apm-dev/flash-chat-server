package domain

type UserRepository interface {
	Add(u User) (string, error)
	FindByUsername(uname string) (*User, error)
	List() ([]User, error)
}
