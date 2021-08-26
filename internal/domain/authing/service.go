package authing

import (
	"fmt"

	"github.com/apm-dev/flash-chat/internal/domain"
	"github.com/apm-dev/flash-chat/pkg/logger"
	"github.com/pkg/errors"
)

type Service interface {
	Register(name, email, pass string) (string, error)
	Login(email, pass string) (string, error)
	Authorize(token string) (*domain.User, error)
}

func NewService(rp domain.UserRepository, jwt *JWTManager) Service {
	return &service{rp, jwt}
}

type service struct {
	repo domain.UserRepository
	jwt  *JWTManager
}

func (svc *service) Register(name, uname, pass string) (string, error) {
	const op string = "domain.authing.service.Register"
	// create domain user object
	user, err := domain.NewUser(name, uname, pass)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}
	// persist user
	_, err = svc.repo.Add(*user)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return "", domain.ErrUserAlreadyExists
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	// generate jwt token with user claims
	token, err := svc.jwt.Generate(*user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	logger.Log(logger.INFO, fmt.Sprintf("%s @%s registered", name, uname))

	return token, nil
}

func (svc *service) Login(email, pass string) (string, error) {
	const op string = "domain.authing.service.Login"
	// fetch user from db
	user, err := svc.repo.FindByUsername(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrUserNotFound
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}
	// check password
	if !user.CheckPassword(pass) {
		logger.Log(logger.INFO, errors.Wrap(domain.ErrWrongCredentials, user.Username).Error())
		return "", domain.ErrWrongCredentials
	}
	// generate jwt token with user claims
	token, err := svc.jwt.Generate(*user)
	if err != nil {
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return "", domain.ErrInternalServer
	}

	logger.Log(logger.INFO, fmt.Sprintf(
		"%s @%s logged-in", user.Name, user.Username,
	))

	return token, nil
}

func (svc *service) Authorize(token string) (*domain.User, error) {
	const op string = "domain.authing.service.Authorize"

	// verify and get claims of token
	claims, err := svc.jwt.Verify(token)
	// handle error
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) {
			logger.Log(logger.INFO, errors.Wrap(err, op).Error())
			return nil, domain.ErrInvalidToken
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	user, err := svc.repo.FindByUsername(claims.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		logger.Log(logger.ERROR, errors.Wrap(err, op).Error())
		return nil, domain.ErrInternalServer
	}

	return user, nil
}
