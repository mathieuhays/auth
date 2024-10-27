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

type registerTemplate interface {
	Register(writer io.Writer, form *forms.Form) error
}

func emailFieldValidation(field *forms.Field, form *forms.Form) {
	if field.Value == "" {
		field.Error = fmt.Errorf("this field is required")
		return
	}

	if validate.Email(field.Value) != nil {
		field.Error = fmt.Errorf("invalid email")
		return
	}
}

func passwordFieldValidation(field *forms.Field, form *forms.Form) {
	if field.Value == "" {
		field.Error = fmt.Errorf("this field is required")
		return
	}

	err := validate.Password(field.Value)
	if err != nil {
		field.Error = err
	}
}

func RegisterHandler(tpl registerTemplate, userService user.ServiceInterface) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		emailField := forms.Field{
			Name:     "email",
			Validate: emailFieldValidation,
		}
		emailConfirmField := forms.Field{
			Name:     "email_confirm",
			Validate: emailFieldValidation,
		}
		passwordField := forms.Field{
			Name:     "password",
			Validate: passwordFieldValidation,
		}
		passwordConfirmField := forms.Field{
			Name:     "password_confirm",
			Validate: passwordFieldValidation,
		}

		registrationForm := forms.NewForm(emailField, emailConfirmField, passwordField, passwordConfirmField)
		registrationForm.SetValidation(func(form *forms.Form) {
			hasEmailErrors := form.Fields["email"].Error == nil && form.Fields["email_confirm"] == nil
			hasPasswordErrors := form.Fields["password"].Error == nil && form.Fields["password_confirm"] == nil

			if !hasEmailErrors && form.Fields["email"].Value != form.Fields["email_confirm"].Value {
				form.Fields["email"].Error = fmt.Errorf("emails do not match")
				form.Fields["email_confirm"].Error = fmt.Errorf("")
			}

			if !hasPasswordErrors && form.Fields["password"].Value != form.Fields["password_confirm"].Value {
				form.Fields["password"].Error = fmt.Errorf("passwords do not match")
			}
		})

		if request.Method == http.MethodPost {
			registrationForm.LoadValuesFromRequest(request)
			registrationForm.Validate()

			if !registrationForm.HasErrors() {
				u, err := userService.Register(
					registrationForm.Fields["email"].Value,
					registrationForm.Fields["password"].Value)
				if err == nil {
					log.Printf("registration successful")
					_, session, err2 := userService.Login(u)
					if err2 == nil {
						err3 := userService.SetAuthResponse(writer, session)
						if err3 != nil {
							log.Printf("register: failed setAuthResponse: %s", err3)
						}
						http.Redirect(writer, request, "/dashboard", http.StatusFound)
						return
					}

					log.Printf("register: login error: %s", err2)
					http.Redirect(writer, request, "/login", http.StatusFound)
					return
				}

				registrationForm.Error = fmt.Errorf("error while registering: %s", err)
			}
		}

		if err := tpl.Register(writer, registrationForm); err != nil {
			log.Printf("login error: %s", err)
		}
	})
}
