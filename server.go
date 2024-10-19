package auth

import "net/http"

func NewServer(tpl TemplateEngineInterface) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// 1. home
	// 2. dashboard -- use requireLogin middleware
	// 3. register
	// 4. login
	// 5. lost password

	return mux
}
