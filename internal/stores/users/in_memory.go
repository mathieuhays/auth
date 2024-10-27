package users

import (
	"github.com/google/uuid"
	"github.com/mathieuhays/auth/internal/validate"
	"sync"
	"time"
)

type UserMemoryStore struct {
	items map[uuid.UUID]User
	mu    sync.RWMutex
}

func NewUserMemoryStore() *UserMemoryStore {
	return &UserMemoryStore{items: make(map[uuid.UUID]User), mu: sync.RWMutex{}}
}

func (u *UserMemoryStore) Create(user User) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	emptyUUID := uuid.UUID{}
	if user.ID == emptyUUID {
		user.ID = uuid.New()
	}

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
	u.mu.RLock()
	defer u.mu.RUnlock()

	if localUser, ok := u.items[id]; ok {
		return &localUser, nil
	}

	return nil, ErrUserNotFound
}

func (u *UserMemoryStore) GetByEmail(email string) (*User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	if validate.Email(email) == nil {
		for _, user := range u.items {
			if user.Email == email {
				return &user, nil
			}
		}
	}

	return nil, ErrUserNotFound
}

func (u *UserMemoryStore) Update(user User) (*User, error) {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.items[user.ID]; !ok {
		return nil, ErrUserNotFound
	}

	u.items[user.ID] = user
	localUser := u.items[user.ID]

	return &localUser, nil
}

func (u *UserMemoryStore) Delete(id uuid.UUID) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	delete(u.items, id)
	return nil
}
