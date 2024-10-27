package handlers

import (
	"github.com/mathieuhays/auth/internal/services/user"
	"github.com/mathieuhays/auth/internal/stores/sessions"
	"github.com/mathieuhays/auth/internal/stores/users"
	"io"
	"log"
	"net/http"
)

type dashboardTemplates interface {
	Dashboard(writer io.Writer, u *users.User, s *sessions.Session) error
}

func DashboardHandler(tpl dashboardTemplates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, session, err := user.RetrieveAuthDetails(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if err := tpl.Dashboard(w, u, session); err != nil {
			log.Printf("template error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
