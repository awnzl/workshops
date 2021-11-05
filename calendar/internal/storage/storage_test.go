package storage

import (
	"log"
	"os"
	"testing"

	"github.com/matryer/is"
)

var logger = log.New(os.Stdout, "Log: ", log.Lshortfile)

func TestConnect(t *testing.T) {
	is := is.New(t)

	db := NewPostgresQL(logger)
	err := db.Connect("postgres://gouser:gopassword@localhost:5432/gotest?sslmode=disable&connect_timeout=20")
	is.NoErr(err)

	db.Close()
}
