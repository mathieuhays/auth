package templates

import (
	"html/template"
	"io"
)

type Engine struct {
	tpl *template.Template
}

func NewEngine(tpl *template.Template) Engine {
	return Engine{tpl: tpl}
}

func (t Engine) Index(writer io.Writer) error {
	return t.tpl.ExecuteTemplate(writer, "index", nil)
}

func (t Engine) Error(writer io.Writer, title, description string) error {
	return t.tpl.ExecuteTemplate(writer, "error", struct {
		Error struct {
			Title       string
			Description string
		}
	}{
		Error: struct {
			Title       string
			Description string
		}{Title: title, Description: description},
	})
}
