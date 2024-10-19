package users

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

type UserStoreInterface interface {
	Create(user User) (*User, error)
	Get(id uuid.UUID) (*User, error)
	Update(user User) (*User, error)
	Delete(id uuid.UUID) error
}
