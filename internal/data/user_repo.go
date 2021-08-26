package data

import (
	"github.com/apm-dev/flash-chat/internal/domain"
	"github.com/pkg/errors"
)

type UserDataSource interface {
	Insert(domain.User) (string, error)
	GetByUsername(string) (*domain.User, error)
	GetAll() ([]domain.User, error)
}

func NewUserRepo(ds UserDataSource) domain.UserRepository {
	return &userRepo{ds}
}

type userRepo struct {
	ds UserDataSource
}

func (r *userRepo) Add(u domain.User) (string, error) {
	const op string = "data.repository.user_repo.Add"

	id, err := r.ds.Insert(u)
	if err != nil {
		return "", errors.Wrap(err, op)
	}
	return id, nil
}

func (r *userRepo) FindByUsername(uname string) (*domain.User, error) {
	const op string = "data.repository.user_repo.FindByEmail"

	user, err := r.ds.GetByUsername(uname)
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return user, nil
}

func (r *userRepo) List() ([]domain.User, error) {
	const op string = "data.repository.user_repo.FindAll"

	users, err := r.ds.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, op)
	}

	return users, nil
}
