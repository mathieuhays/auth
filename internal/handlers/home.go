package handlers

import (
	"github.com/mathieuhays/auth"
	"net/http"
)

func HomeHandler(tpl auth.TemplateEngineInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tpl.Index(w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
