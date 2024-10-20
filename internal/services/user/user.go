package user

import (
	"errors"
	"github.com/mathieuhays/auth/internal/stores/sessions"
	"github.com/mathieuhays/auth/internal/stores/users"
)

type ServiceInterface interface {
	Login(email, password string) (*users.User, *sessions.Session, error)
	LoginWithToken(sessionToken string) (*users.User, *sessions.Session, error)
	Register(email, password string) (*users.User, error)
}

type Service struct {
	userStore    users.UserStoreInterface
	sessionStore sessions.SessionStoreInterface
}

func NewService(userStore users.UserStoreInterface, sessionStore sessions.SessionStoreInterface) *Service {
	return &Service{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (s Service) Login(email, password string) (*users.User, *sessions.Session, error) {
	return nil, nil, errors.New("not implemented")
}

func (s Service) LoginWithToken(sessionToken string) (*users.User, *sessions.Session, error) {
	return nil, nil, errors.New("not implemented")
}

func (s Service) Register(email, password string) (*users.User, error) {
	return nil, errors.New("not implemented")
}
