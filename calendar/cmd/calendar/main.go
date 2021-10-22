package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"calendar/internal/api"
	"calendar/internal/app"
	"calendar/internal/middleware"
)

/*
- Start from utilizing standard packages: net/http; !
- Use net/http/httptest for creation tests;
- Find out how to convert & store time between different timezones with package time;
- All data should be stored in memory;
- Use Postman to validate your server;
- Use is and moq in tests;
*/

const (
	port = 8000
)

func main() {
	logger := log.New(os.Stdout, "Log: ", log.Lshortfile)

	router := mux.NewRouter()
	// authenticator := authapp.New()
	application := app.New(logger)

	handlers := api.New(application, logger)
	handlers.RegisterHandlers(
		router,
		middleware.Logger(logger),
	)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	logger.Println("start listening port", port)

	if err := s.ListenAndServe(); err != nil {
		logger.Panic("server error", err)
	}
}
