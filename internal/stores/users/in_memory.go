package users

import (
	"github.com/google/uuid"
	"time"
)

type UserMemoryStore struct {
	items map[uuid.UUID]User
}

func NewUserMemoryStore() *UserMemoryStore {
	return &UserMemoryStore{items: make(map[uuid.UUID]User)}
}

func (u *UserMemoryStore) Create(user User) (*User, error) {
	if _, ok := u.items[user.ID]; ok {
		return nil, ErrUserAlreadyExist
	}

	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now().UTC()
	}

	u.items[user.ID] = user
	localUser := u.items[user.ID]

	return &localUser, nil
}

func (u *UserMemoryStore) Get(id uuid.UUID) (*User, error) {
	if localUser, ok := u.items[id]; ok {
		return &localUser, nil
	}

	return nil, ErrUserNotFound
}

func (u *UserMemoryStore) Update(user User) (*User, error) {
	if _, ok := u.items[user.ID]; !ok {
		return nil, ErrUserNotFound
	}

	u.items[user.ID] = user
	localUser := u.items[user.ID]

	return &localUser, nil
}

func (u *UserMemoryStore) Delete(id uuid.UUID) error {
	delete(u.items, id)
	return nil
}
