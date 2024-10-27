package sessions

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type SessionMemoryStore struct {
	items map[uuid.UUID]Session
	mu    sync.RWMutex
}

func NewSessionMemoryStore() *SessionMemoryStore {
	return &SessionMemoryStore{make(map[uuid.UUID]Session), sync.RWMutex{}}
}

func (s *SessionMemoryStore) Create(session Session) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	emptyUUID := uuid.UUID{}

	if session.ID == emptyUUID {
		session.ID = uuid.New()
	}

	if _, ok := s.items[session.ID]; ok {
		return nil, ErrSessionAlreadyExist
	}

	if session.UserID == emptyUUID {
		return nil, ErrSessionMissingUserID
	}

	if session.CreatedAt.IsZero() {
		session.CreatedAt = time.Now().UTC()
	}

	s.items[session.ID] = session
	localSession := s.items[session.ID]

	return &localSession, nil
}

func (s *SessionMemoryStore) Get(id uuid.UUID) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if session, ok := s.items[id]; ok {
		return &session, nil
	}

	return nil, ErrSessionNotFound
}

func (s *SessionMemoryStore) GetForUser(userID uuid.UUID) ([]Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var sessions []Session

	for _, session := range s.items {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (s *SessionMemoryStore) GetForToken(token string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, session := range s.items {
		if session.Token == token {
			return &session, nil
		}
	}

	return nil, ErrSessionNotFound
}

func (s *SessionMemoryStore) Update(session Session) (*Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[session.ID]; !ok {
		return nil, ErrSessionNotFound
	}

	s.items[session.ID] = session
	localSession := s.items[session.ID]

	return &localSession, nil
}

func (s *SessionMemoryStore) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, id)
	return nil
}
