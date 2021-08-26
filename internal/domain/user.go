package domain

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func NewUser(name, email, pass string) (*User, error) {
	const op string = "domain.user.NewUser"

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return &User{
		Name:     name,
		Username: email,
		Password: string(hash),
	}, nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (user *User) Clone() *User {
	return &User{
		Name:     user.Name,
		Username: user.Username,
		Password: user.Password,
	}
}
