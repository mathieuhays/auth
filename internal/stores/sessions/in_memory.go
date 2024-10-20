package sessions

import (
	"github.com/google/uuid"
	"time"
)

type SessionMemoryStore struct {
	items map[uuid.UUID]Session
}

func NewSessionMemoryStore() *SessionMemoryStore {
	return &SessionMemoryStore{make(map[uuid.UUID]Session)}
}

func (s *SessionMemoryStore) Create(session Session) (*Session, error) {
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
	if session, ok := s.items[id]; ok {
		return &session, nil
	}

	return nil, ErrSessionNotFound
}

func (s *SessionMemoryStore) GetForUser(userID uuid.UUID) ([]Session, error) {
	var sessions []Session

	for _, session := range s.items {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}

	return sessions, nil
}

func (s *SessionMemoryStore) Update(session Session) (*Session, error) {
	if _, ok := s.items[session.ID]; !ok {
		return nil, ErrSessionNotFound
	}

	s.items[session.ID] = session
	localSession := s.items[session.ID]

	return &localSession, nil
}

func (s *SessionMemoryStore) Delete(id uuid.UUID) error {
	delete(s.items, id)
	return nil
}
