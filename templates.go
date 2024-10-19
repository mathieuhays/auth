package auth

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/fragments templates
var templates embed.FS

type TemplateEngineInterface interface {
	Index(writer io.Writer) error
}

type TemplateEngine struct {
	tpl *template.Template
}

func NewTemplateEngine() (*TemplateEngine, error) {
	tpl, err := template.New("auth").Funcs(template.FuncMap{
		"assets": func(file string) string {
			return "/static/" + file
		},
	}).ParseFS(templates, "templates/*.gohtml", "templates/fragments/*.gohtml")
	if err != nil {
		return nil, err
	}

	return &TemplateEngine{tpl: tpl}, nil
}

func (t TemplateEngine) Index(writer io.Writer) error {
	return t.tpl.ExecuteTemplate(writer, "index", nil)
}
