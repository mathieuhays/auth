package handlers

import (
	"io"
	"log"
	"net/http"
)

type errorTemplates interface {
	Error(writer io.Writer, title, description string) error
}

func ErrorHandler(tpl errorTemplates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		if err := tpl.Error(w, "Error 404", "Page not found"); err != nil {
			log.Printf("template error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}
