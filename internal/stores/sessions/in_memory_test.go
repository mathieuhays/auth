package sessions

import (
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestSessionMemoryStore_Create(t *testing.T) {
	store := NewSessionMemoryStore()

	t.Run("ensure ID gets set", func(t *testing.T) {
		session := Session{
			ID:     uuid.UUID{},
			UserID: uuid.New(),
		}

		s, err := store.Create(session)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		emptyUUID := uuid.UUID{}
		if s.ID == emptyUUID {
			t.Fatalf("session ID is empty. should be set if not provided.")
		}
	})

	t.Run("require user id", func(t *testing.T) {
		session := Session{
			ID:        uuid.New(),
			UserID:    uuid.UUID{},
			Token:     "test",
			CSRFToken: "test",
		}

		_, err := store.Create(session)
		if err == nil {
			t.Fatalf("error expected but none returned")
		}

		if !errors.Is(err, ErrSessionMissingUserID) {
			t.Errorf("unexpected error returned. expected: %s. got: %s", ErrSessionMissingUserID, err)
		}
	})

	t.Run("ensure CreatedAt gets set", func(t *testing.T) {
		session := Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}

		s, err := store.Create(session)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if s.CreatedAt.IsZero() {
			t.Fatalf("CreatedAt did not get a default value")
		}
	})

	t.Run("ensure CreatedAt doest not get overridden", func(t *testing.T) {
		session := Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}

		s, err := store.Create(session)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if s.CreatedAt.IsZero() {
			t.Fatalf("CreatedAt did not get a default value")
		}
	})
}

func TestSessionMemoryStore_Get(t *testing.T) {
	store := NewSessionMemoryStore()
	existingSession := Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "test",
		CSRFToken: "test",
		CreatedAt: time.Now().UTC(),
	}
	_, err := store.Create(existingSession)
	if err != nil {
		t.Fatalf("unexpected error when creating sessions: %s", err)
	}

	s, err := store.Get(existingSession.ID)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if s.ID != existingSession.ID {
		t.Fatalf("unexpected element. expected: %s. got: %s. obj: %v", existingSession.ID, s.ID, s)
	}
}

func TestSessionMemoryStore_GetForUser(t *testing.T) {
	store := NewSessionMemoryStore()
	userID := uuid.New()

	session1 := Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "test",
		CSRFToken: "test",
		CreatedAt: time.Now().UTC(),
	}
	_, err := store.Create(session1)
	if err != nil {
		t.Fatalf("unexpected error when creating session: %s", err)
	}

	session2 := Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "test",
		CSRFToken: "test",
		CreatedAt: time.Now().UTC(),
	}
	_, err = store.Create(session2)
	if err != nil {
		t.Fatalf("unexpected error when creating session: %s", err)
	}

	session3 := Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "test",
		CSRFToken: "test",
		CreatedAt: time.Now().UTC(),
	}
	_, err = store.Create(session3)
	if err != nil {
		t.Fatalf("unexpected error when creating session: %s", err)
	}

	unknownUserSessions, err := store.GetForUser(uuid.New())
	if err != nil {
		t.Errorf("unexpected error when retrieving sessions for unknown user: %s", err)
	} else if len(unknownUserSessions) != 0 {
		t.Errorf("unexpected amount of sessions found for unknown user. expected: 0. got: %d", len(unknownUserSessions))
	}

	userSessions, err := store.GetForUser(userID)
	if err != nil {
		t.Errorf("unexpected error when retrieving sessions for known user: %s", err)
	} else if len(userSessions) != 2 {
		t.Errorf("unexpected amount of sessions found for known user. expected 2. got: %d", len(userSessions))
	}
}

func TestSessionMemoryStore_GetForToken(t *testing.T) {
	store := NewSessionMemoryStore()
	token := "test_token"
	session := Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     token,
		CSRFToken: "csrf_token",
		CreatedAt: time.Now(),
	}

	_, err := store.Create(session)
	if err != nil {
		t.Fatalf("unexpected error when creating session: %s", err)
	}

	s, err := store.GetForToken(token)
	if err != nil {
		t.Fatalf("unexpected error while retrieving session: %s", err)
	}

	if s.ID != session.ID {
		t.Errorf("session does not match. expected: %s. got: %s", session.ID, s.ID)
	}
}

func TestSessionMemoryStore_Update(t *testing.T) {
	store := NewSessionMemoryStore()
	session := Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "test",
		CSRFToken: "test",
		CreatedAt: time.Now(),
	}
	_, err := store.Create(session)
	if err != nil {
		t.Fatalf("unexpected error when creating session: %s", err)
	}

	unknownSession := Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "",
		CSRFToken: "",
		CreatedAt: time.Now(),
	}
	_, err = store.Update(unknownSession)
	if err == nil {
		t.Errorf("no error thrown when expected for unknown session")
	} else if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("unexpected error for unknown session: %s", err)
	}

	newSession := session
	newSession.Token = "updated"
	updatedSession, err := store.Update(newSession)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if updatedSession.Token != newSession.Token {
		t.Fatalf("value is not updated. expected: %s. got: %s", newSession.Token, updatedSession.Token)
	}
}

func TestSessionMemoryStore_Delete(t *testing.T) {
	store := NewSessionMemoryStore()
	session := Session{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Token:     "",
		CSRFToken: "",
		CreatedAt: time.Now(),
	}

	if _, err := store.Create(session); err != nil {
		t.Fatalf("unexpected error when creating session: %s", err)
	}

	if err := store.Delete(session.ID); err != nil {
		t.Fatalf("unexpected error when deleting session: %s", err)
	}

	_, err := store.Get(session.ID)
	if err == nil {
		t.Fatalf("error expected after retrieving the session we just deleted but none returned")
	}

	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("unexpected error when retrieving deleted session: %s", err)
	}
}
