package stores

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type User struct {
	ID             uuid.UUID
	Email          string
	EmailConfirmed *time.Time
	PasswordHash   string
}

type UserStore struct {
	items map[uuid.UUID]User
}

type UserStoreInterface interface {
	Create(user User) (*User, error)
	Get(id uuid.UUID) (*User, error)
	Update(user User) (*User, error)
	Delete(id uuid.UUID) error
}

func NewUserStore() *UserStore {
	return &UserStore{items: make(map[uuid.UUID]User)}
}

func (u *UserStore) Create(user User) (*User, error) {
	if _, ok := u.items[user.ID]; ok {
		return nil, ErrUserAlreadyExist
	}

	u.items[user.ID] = user
	localUser := u.items[user.ID]

	return &localUser, nil
}

func (u *UserStore) Get(id uuid.UUID) (*User, error) {
	if localUser, ok := u.items[id]; ok {
		return &localUser, nil
	}

	return nil, ErrUserNotFound
}

func (u *UserStore) Update(user User) (*User, error) {
	if _, ok := u.items[user.ID]; !ok {
		return nil, ErrUserNotFound
	}

	u.items[user.ID] = user
	localUser := u.items[user.ID]

	return &localUser, nil
}

func (u *UserStore) Delete(id uuid.UUID) error {
	delete(u.items, id)
	return nil
}
