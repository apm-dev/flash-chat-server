package listing

import (
	"github.com/apm-dev/flash-chat/internal/domain"
	"github.com/apm-dev/flash-chat/pkg/logger"
	"github.com/pkg/errors"
)

type Service interface {
	Users() ([]domain.User, error)
}

type service struct {
	repo domain.UserRepository
}

func NewService(rp domain.UserRepository) Service {
	return &service{rp}
}

func (s *service) Users() ([]domain.User, error) {
	const op string = "domain.listing.Users"

	users, err := s.repo.List()
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	if len(users) == 0 {
		return nil, domain.ErrUserNotFound
	}

	return users, nil
}
