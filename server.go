package auth

import (
	"github.com/mathieuhays/auth/internal/handlers"
	"github.com/mathieuhays/auth/internal/services/user"
	"github.com/mathieuhays/auth/internal/templates"
	"log"
	"net/http"
)

func NewServer(tpl *templates.Engine, userService user.ServiceInterface) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /", handlers.ErrorHandler(tpl))
	mux.Handle("GET /{$}", handlers.HomeHandler(tpl))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	mux.Handle("/login", handlers.LoginHandler(userService))

	// 1. home
	// 2. dashboard -- use requireLogin middleware
	// 3. register
	// 4. login
	// 5. lost password

	return loggerMiddleware(mux)
}

/*
ANSI Color codes, 30-base, 90 range for bright variant
0: Black
1: Red
2: Green
3: Yellow
4: Blue
5: Magenta
6: Cyan
7: White
*/
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("\u001b[32m%s\u001b[0m %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
