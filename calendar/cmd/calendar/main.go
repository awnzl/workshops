package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"calendar/internal/api"
	"calendar/internal/app"
	"calendar/internal/auth"
	"calendar/internal/storage"
)

/*
- Start from utilizing standard packages: net/http; !
- Use net/http/httptest for creation tests;
- Find out how to convert & store time between different timezones with package time;
- All data should be stored in memory;
- Use Postman to validate your server;
- Use is and moq in tests;
*/

// debug purpose
const (
	port            = 8000
	secretKey       = "Pr3ttyS3cr3tK3y"
	tokenExpiration = 24
	DbURL           = "postgres://gouser:gopassword@localhost:5432/gotest?sslmode=disable"
)

func serve(ctx context.Context, logger *log.Logger) error {
	db := storage.NewPostgresQL(logger)
	if err := db.Connect(dbURL); err != nil {
		logger.Fatal("database connecting fail:", err)
	}
	defer db.Close()

	authApp := auth.New(secretKey, tokenExpiration)
	application := app.New(logger, authApp, db)

	router := mux.NewRouter()
	handlers := api.New(application, authApp, logger)
	handlers.RegisterHandlers(
		router,
	)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", err)
		}
	}()

	logger.Println("start listening port", port)

	<-ctx.Done()

	logger.Println("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.Shutdown(ctxShutDown)
	if err != nil {
		logger.Fatal("server ShutDown failed:", err)
	}

	logger.Printf("server exited properly")

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	logger := log.New(os.Stdout, "Log: ", log.Lshortfile)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		logger.Println("system call", <-ch)
		cancel()
	}()

	if err := serve(ctx, logger); err != nil {
		logger.Println("failed to serve:", err)
	}
}
