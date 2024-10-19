package handlers

import (
	"github.com/mathieuhays/auth/internal/asserts"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorHandlerTpl struct {
	callback func(title, description string) error
}

func (e errorHandlerTpl) Error(writer io.Writer, title, description string) error {
	return e.callback(title, description)
}

func TestErrorHandler(t *testing.T) {
	var responseTitle string
	var responseDescription string

	serverHandler := ErrorHandler(errorHandlerTpl{
		callback: func(title, description string) error {
			responseTitle = title
			responseDescription = description
			return nil
		},
	})
	req, err := http.NewRequest(http.MethodGet, "/123", nil)
	if err != nil {
		t.Fatalf("request creation error: %s", err)
	}

	response := httptest.NewRecorder()
	serverHandler.ServeHTTP(response, req)

	asserts.StatusCode(t, response, http.StatusNotFound)

	if !strings.Contains(responseTitle, "404") {
		t.Errorf("title does not contain 404. got: %s", responseTitle)
	}

	if !strings.Contains(responseDescription, "not found") {
		t.Errorf("description does not contain 'not found'. got: %s", responseDescription)
	}
}
