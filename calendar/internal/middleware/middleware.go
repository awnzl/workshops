package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Logger(l *log.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Println("Request", fmt.Sprintln("URI", r.RequestURI), fmt.Sprintln("Addr", r.RemoteAddr))
			handler.ServeHTTP(w, r)
		})
	}
}
