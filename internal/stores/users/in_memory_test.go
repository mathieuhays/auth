package users

import (
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestUserMemoryStore_Create(t *testing.T) {
	t.Run("ID fallback", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			Email:     "email",
			CreatedAt: time.Now(),
		}
		u, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		emptyUUID := uuid.UUID{}
		if u.ID == emptyUUID {
			t.Fatalf("failed to generate a fallback ID")
		}
	})

	t.Run("ID fallback no override", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID: uuid.New(),
		}

		u, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if u.ID != user.ID {
			t.Fatalf("explicit ID has changed")
		}
	})

	t.Run("CreadtedAt fallback", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:    uuid.New(),
			Email: "test",
		}

		u, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if u.CreatedAt.IsZero() {
			t.Fatalf("CreatedAt is empty. Fallback value should have been applied")
		}
	})

	t.Run("CreatedAt fallback override", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:        uuid.New(),
			CreatedAt: time.Now().Add(time.Hour * -2),
		}

		u, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if !u.CreatedAt.Equal(user.CreatedAt) {
			t.Fatalf("unexpected CreatedAt value. expected: %v. got: %v", user.CreatedAt, u.CreatedAt)
		}
	})

	t.Run("success", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:             uuid.New(),
			Email:          "test",
			EmailConfirmed: nil,
			PasswordHash:   "test",
		}

		_, err := store.Create(user)
		if err != nil {
			t.Fatalf("Create: unexpected error: %s", err)
		}

		u, err := store.Get(user.ID)
		if err != nil {
			t.Fatalf("Get: unexpected error: %s", err)
		}

		if u.Email != user.Email {
			t.Fatalf("retrieved entry does not matched the one created. expected: %s. got: %s", user.Email, u.Email)
		}
	})
}

func TestUserMemoryStore_Get(t *testing.T) {
	t.Run("fetch existing record", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:    uuid.New(),
			Email: "test",
		}

		_, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error while creating record: %s", err)
		}

		u, err := store.Get(user.ID)
		if err != nil {
			t.Fatalf("unexpected error while getting user: %s", err)
		}

		if u.Email != user.Email {
			t.Fatalf("records do not match. expected: %s. got: %s", user.Email, u.Email)
		}
	})

	t.Run("fetch missing record", func(t *testing.T) {
		store := NewUserMemoryStore()
		_, err := store.Get(uuid.New())
		if err == nil {
			t.Fatalf("did not get an error while fetching a missing record")
		}

		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("unexpected error. expected: %s. got: %s", ErrUserNotFound, err)
		}
	})
}

func TestUserMemoryStore_GetByEmail(t *testing.T) {
	t.Run("existing email", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:    uuid.New(),
			Email: "test@example.com",
		}
		_, err := store.Create(user)
		if err != nil {
			t.Fatalf("error while creating user: %s", err)
		}

		u, err := store.GetByEmail(user.Email)
		if err != nil {
			t.Fatalf("unexpected error while retrieving user by email: %s", err)
		}

		if u.ID != user.ID {
			t.Fatalf("IDs do not match. expected: %v. got: %v", user.ID, u.ID)
		}
	})

	t.Run("missing email", func(t *testing.T) {
		store := NewUserMemoryStore()
		_, err := store.GetByEmail("test@example.com")
		if err == nil {
			t.Fatalf("did not get an error")
		}

		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("unexpected error. expected: %s. got: %s", ErrUserNotFound, err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		store := NewUserMemoryStore()
		_, err := store.GetByEmail("invalid")
		if err == nil {
			t.Fatalf("did not get an error")
		}

		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("unexpected error. expected: %s. got: %s", ErrUserNotFound, err)
		}
	})
}

func TestUserMemoryStore_Update(t *testing.T) {
	t.Run("existing record", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:    uuid.New(),
			Email: "test",
		}
		u, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error while creating user: %s", err)
		}

		newUserData := *u
		newUserData.Email = "another@email.com"
		updatedUser, err := store.Update(newUserData)
		if err != nil {
			t.Fatalf("unexpected error while updating user: %s", err)
		}

		if updatedUser.Email == user.Email {
			t.Fatalf("user has not been updated. expected: %s. got: %s", user.Email, updatedUser.Email)
		}
	})

	t.Run("missing record", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:        uuid.New(),
			Email:     "test",
			CreatedAt: time.Now(),
		}

		_, err := store.Update(user)
		if err == nil {
			t.Fatalf("did not get an error")
		}

		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("unexpected error. expected: %s. got: %s", ErrUserNotFound, err)
		}
	})
}

func TestUserMemoryStore_Delete(t *testing.T) {
	t.Run("existing record", func(t *testing.T) {
		store := NewUserMemoryStore()
		user := User{
			ID:    uuid.New(),
			Email: "test",
		}
		_, err := store.Create(user)
		if err != nil {
			t.Fatalf("unexpected error while creating user: %s", err)
		}

		err = store.Delete(user.ID)
		if err != nil {
			t.Fatalf("unexpected error while deleting user: %s", err)
		}

		_, err = store.Get(user.ID)
		if err == nil {
			t.Fatalf("not error returned when retrieving deleted user")
		}

		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("unexpected error when retrieving deleted user. expected: %s. got: %s", ErrUserNotFound, err)
		}
	})

	t.Run("missing record", func(t *testing.T) {
		store := NewUserMemoryStore()

		err := store.Delete(uuid.New())
		if err != nil {
			t.Fatalf("unexpected error returned. should fail silently. got: %s", err)
		}
	})
}
