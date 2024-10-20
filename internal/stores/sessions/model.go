package sessions

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExist  = errors.New("session already exist")
	ErrSessionMissingUserID = errors.New("user ID missing")
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	CSRFToken string
	CreatedAt time.Time
}

type SessionStoreInterface interface {
	Create(session Session) (*Session, error)
	Get(id uuid.UUID) (*Session, error)
	GetForUser(userID uuid.UUID) ([]Session, error)
	Update(session Session) (*Session, error)
	Delete(id uuid.UUID) error
}
