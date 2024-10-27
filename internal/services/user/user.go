package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mathieuhays/auth/internal/stores/sessions"
	"github.com/mathieuhays/auth/internal/stores/users"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

const authCookie = "session_token"

type ContextKey string

const UserContextKey = "user"
const SessionContextKey = "session"

type ServiceInterface interface {
	Login(user *users.User) (*users.User, *sessions.Session, error)
	LoginWithCredentials(email, password string) (*users.User, *sessions.Session, error)
	LoginWithToken(sessionToken string) (*users.User, *sessions.Session, error)
	Register(email, password string) (*users.User, error)
	SetAuthResponse(writer http.ResponseWriter, session *sessions.Session) error
	RetrieveAuthFromRequest(request *http.Request) (*users.User, *sessions.Session, error)
}

type Service struct {
	userStore    users.UserStoreInterface
	sessionStore sessions.SessionStoreInterface
}

func NewService(userStore users.UserStoreInterface, sessionStore sessions.SessionStoreInterface) *Service {
	return &Service{
		userStore:    userStore,
		sessionStore: sessionStore,
	}
}

func (s Service) Login(user *users.User) (*users.User, *sessions.Session, error) {
	session, err := sessions.NewSession(user.ID)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.sessionStore.Create(*session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func (s Service) LoginWithCredentials(email, password string) (*users.User, *sessions.Session, error) {
	user, err := s.userStore.GetByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, err
	}

	return s.Login(user)
}

func (s Service) LoginWithToken(sessionToken string) (*users.User, *sessions.Session, error) {
	session, err := s.sessionStore.GetForToken(sessionToken)
	if err != nil {
		return nil, nil, err
	}

	user, err := s.userStore.Get(session.UserID)
	if err != nil {
		return nil, nil, err
	}

	session.LastUsed = time.Now()
	_, err = s.sessionStore.Update(*session)
	if err != nil {
		log.Printf("session failed to update lastUsed field: %s", err)
	}

	return user, session, nil
}

func (s Service) Register(email, password string) (*users.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := users.User{
		ID:             uuid.New(),
		Email:          email,
		EmailConfirmed: nil,
		PasswordHash:   string(passwordHash),
		CreatedAt:      time.Now(),
	}

	return s.userStore.Create(user)
}

func (s Service) SetAuthResponse(writer http.ResponseWriter, session *sessions.Session) error {
	if session == nil {
		return errors.New("invalid session")
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     authCookie,
		Value:    session.Token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

func (s Service) RetrieveAuthFromRequest(request *http.Request) (*users.User, *sessions.Session, error) {
	cookie, err := request.Cookie(authCookie)
	if err != nil {
		return nil, nil, err
	}

	user, session, err := s.LoginWithToken(cookie.Value)
	if err != nil {
		return nil, nil, err
	}

	session.LastUsed = time.Now()
	_, err = s.sessionStore.Update(*session)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func AugmentRequestWithAuth(request *http.Request, user *users.User, session *sessions.Session) *http.Request {
	ctx := context.WithValue(request.Context(), UserContextKey, *user)
	ctx = context.WithValue(ctx, SessionContextKey, *session)

	return request.WithContext(ctx)
}

func RetrieveAuthDetails(request *http.Request) (*users.User, *sessions.Session, error) {
	user, ok := request.Context().Value(UserContextKey).(users.User)
	if !ok {
		return nil, nil, errors.New("no users found")
	}

	session, ok := request.Context().Value(SessionContextKey).(sessions.Session)
	if !ok {
		return nil, nil, errors.New("no session found")
	}

	return &user, &session, nil
}
