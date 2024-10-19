package handlers

import (
	"io"
	"log"
	"net/http"
)

type homeTemplates interface {
	Index(writer io.Writer) error
}

func HomeHandler(tpl homeTemplates) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := tpl.Index(w); err != nil {
			log.Printf("template error: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
