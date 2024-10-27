package templates

import (
	"github.com/mathieuhays/auth/internal/forms"
	"github.com/mathieuhays/auth/internal/stores/sessions"
	"github.com/mathieuhays/auth/internal/stores/users"
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

func (t Engine) Register(writer io.Writer, form *forms.Form) error {
	return t.tpl.ExecuteTemplate(writer, "register", struct {
		Form *forms.Form
	}{
		Form: form,
	})
}

func (t Engine) Login(writer io.Writer, form *forms.Form) error {
	return t.tpl.ExecuteTemplate(writer, "login", struct {
		Form *forms.Form
	}{
		Form: form,
	})
}

func (t Engine) Dashboard(writer io.Writer, u *users.User, s *sessions.Session) error {
	return t.tpl.ExecuteTemplate(writer, "dashboard", struct {
		User    *users.User
		Session *sessions.Session
	}{
		User:    u,
		Session: s,
	})
}
