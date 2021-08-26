package memory

import (
	"sync"

	"github.com/apm-dev/flash-chat/internal/data"
	"github.com/apm-dev/flash-chat/internal/domain"
)

type userDS struct {
	sync.RWMutex
	users map[string]domain.User
}

func NewUserDS() data.UserDataSource {
	return &userDS{
		users: make(map[string]domain.User),
	}
}

func (ds *userDS) Insert(u domain.User) (string, error) {
	ds.Lock()
	defer ds.Unlock()

	if _, ok := ds.users[u.Username]; ok {
		return "", domain.ErrUserAlreadyExists
	}

	ds.users[u.Username] = u
	return u.Username, nil
}

func (ds *userDS) GetByUsername(uname string) (*domain.User, error) {
	ds.RLock()
	defer ds.RUnlock()

	u, ok := ds.users[uname]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	return u.Clone(), nil
}

func (ds *userDS) GetAll() ([]domain.User, error) {
	ds.RLock()
	defer ds.RUnlock()

	if len(ds.users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	users := make([]domain.User, 0, len(ds.users))
	for _, u := range ds.users {
		users = append(users, u)
	}

	return users, nil
}
