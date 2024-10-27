package sessions

import (
	"crypto/rand"
	"encoding/hex"
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
	LastUsed  time.Time
}

func generateToken() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func NewSession(userID uuid.UUID) (*Session, error) {
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	csrfToken, err := generateToken()
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		CSRFToken: csrfToken,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}, nil
}

type SessionStoreInterface interface {
	Create(session Session) (*Session, error)
	Get(id uuid.UUID) (*Session, error)
	GetForUser(userID uuid.UUID) ([]Session, error)
	GetForToken(token string) (*Session, error)
	Update(session Session) (*Session, error)
	Delete(id uuid.UUID) error
}
