package handlers

import (
	"github.com/mathieuhays/auth/internal/services/user"
	"net/http"
)

func LoginHandler(userService user.ServiceInterface) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	})
}
