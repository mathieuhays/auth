package auth

import (
	"embed"
	"html/template"
)

//go:embed templates/fragments templates
var templatesFS embed.FS

func Templates() (*template.Template, error) {
	return template.New("auth").Funcs(template.FuncMap{
		"assets": func(file string) string {
			return "/static/" + file
		},
	}).ParseFS(templatesFS, "templates/*.gohtml", "templates/fragments/*.gohtml")
}
