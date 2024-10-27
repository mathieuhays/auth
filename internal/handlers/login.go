package handlers

import (
	"fmt"
	"github.com/mathieuhays/auth/internal/forms"
	"github.com/mathieuhays/auth/internal/services/user"
	"github.com/mathieuhays/auth/internal/validate"
	"io"
	"log"
	"net/http"
)

type loginTemplate interface {
	Login(writer io.Writer, form *forms.Form) error
}

func LoginHandler(tpl loginTemplate, userService user.ServiceInterface) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		emailField := forms.Field{
			Name: "email",
			Validate: func(field *forms.Field, form *forms.Form) {
				if field.Value == "" {
					field.Error = fmt.Errorf("this field is required")
					return
				}

				if validate.Email(field.Value) != nil {
					field.Error = fmt.Errorf("invalid email")
					return
				}
			},
		}
		passwordField := forms.Field{
			Name: "password",
			Validate: func(field *forms.Field, form *forms.Form) {
				if field.Value == "" {
					field.Error = fmt.Errorf("this field is required")
				}
			},
		}
		loginForm := forms.NewForm(emailField, passwordField)

		if request.Method == http.MethodPost {
			loginForm.LoadValuesFromRequest(request)
			loginForm.Validate()

			if !loginForm.HasErrors() {
				_, s, err := userService.LoginWithCredentials(
					loginForm.Fields["email"].Value,
					loginForm.Fields["password"].Value)
				if err == nil {
					err2 := userService.SetAuthResponse(writer, s)
					if err2 == nil {
						http.Redirect(writer, request, "/dashboard", http.StatusFound)
						return
					}

					log.Printf("login: setAuthResponse error: %s", err2)
					loginForm.Error = fmt.Errorf("something went wrong. please try agin")
				} else {
					log.Printf("login error: %s", err)
					loginForm.Error = fmt.Errorf("invalid credentials")
				}
			}
		}

		if err := tpl.Login(writer, loginForm); err != nil {
			log.Printf("login error: %s", err)
		}
	})
}
