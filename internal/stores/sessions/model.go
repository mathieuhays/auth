package sessions

import (
	"github.com/google/uuid"
	"time"
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
	GetForUser(userID uuid.UUID) (*Session, error)
	Update(session Session) (*Session, error)
	Delete(id uuid.UUID) error
}
