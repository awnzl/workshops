package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"calendar/internal/app"
	"calendar/internal/auth"
	"calendar/internal/helpers"
)

func Logger(l *log.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Println("Request", r.Method, r.RequestURI, fmt.Sprintln("Addr", r.RemoteAddr))
			handler.ServeHTTP(w, r)
		})
	}
}

func Authorization(l *log.Logger, authApp app.Authentication) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "No Authorization header provided", http.StatusForbidden)
				return
			}

			extracted := strings.Split(authHeader, "Bearer ")
			if len(extracted) != 2 {
				http.Error(w, "Incorrect Format of Authorization Token", http.StatusBadRequest)
				return
			}

			clientToken := strings.TrimSpace(extracted[1])

			claims, err := authApp.ValidateJWT(clientToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			c, ok := claims.(*auth.JWTClaim)
			if !ok {
				http.Error(w, "can't extract authentication claims", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), helpers.CtxValKey("username"), c.Username)
			r = r.WithContext(ctx)

			handler.ServeHTTP(w, r)
		})
	}
}
