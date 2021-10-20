package api

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/matryer/is"
)

// test suit
var ts = struct {
	api *API
}{
	api: New(log.New(os.Stdout, "Log: ", log.Lshortfile)),
}

func TestLogin(t *testing.T) {
	is := is.New(t)

	w := httptest.NewRecorder()
	w.Code = 1 // change default status 200

	req := httptest.NewRequest(
		"POST",
		"/login",
		strings.NewReader(`{"Username": "test", "Password": "test"}`),
	)

	ts.api.login(w, req)

	resp := w.Result()
	is.Equal(http.StatusOK, resp.StatusCode)
}

// func TestLogout(t *testing.T) {
// 	is := is.New(t)

// 	w := &httptest.ResponseRecorder{
// 		HeaderMap: make(http.Header),
// 		Body:      new(bytes.Buffer),
// 		Code: 1,
// 	}


// 	req := httptest.NewRequest(
// 		"GET",
// 		"/logout",
// 		nil,
// 	)

// 	ts.api.logout(w, req)

// 	resp := w.Result()
// 	is.Equal(http.StatusOK, resp.StatusCode)
// }
