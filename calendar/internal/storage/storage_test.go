package storage

import (
	"calendar/internal/auth"
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)

const (
	secretKey  = "3xtr3m3lyS3cr3tK3y"
	expiration = 1
)

var logger = log.New(os.Stdout, "Log: ", log.Lshortfile)

func TestConnect(t *testing.T) {
	is := is.New(t)
	authApp := auth.New(secretKey, expiration)

	db := NewPostgresQL(logger, authApp)
	err := db.Connect("postgres://gouser:gopassword@localhost:5432/gotest?sslmode=disable&connect_timeout=20")
	is.NoErr(err)

	db.Close()
}
