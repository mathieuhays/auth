package stores

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID
	Email          string
	EmailConfirmed *time.Time
	Password       string
}

type UserStore struct {
	items map[uuid.UUID]User
}

type UserStoreInterface interface {
	Create(user User) (*User, error)
	Get(id uuid.UUID) (*User, error)
}

func NewUserStore() *UserStore {
	return &UserStore{items: make(map[uuid.UUID]User)}
}

func (u *UserStore) Create(user User) (*User, error) {
	return nil, errors.New("not implemented")
}

func (u *UserStore) Get(id uuid.UUID) (*User, error) {
	return nil, errors.New("not implemented")
}
