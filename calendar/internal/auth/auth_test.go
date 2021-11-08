package auth

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

const (
	secretKey  = "3xtr3m3lyS3cr3tK3y"
	expiration = 1
)

func TestGenerateToken(t *testing.T) {
	authApp := New(secretKey, expiration)

	is := is.New(t)
	generatedToken, err := authApp.GenerateJWT("UserName")
	is.NoErr(err)

	os.Setenv("testToken", generatedToken)
}

func TestValidateToken(t *testing.T) {
	encodedToken := os.Getenv("testToken")
	authApp := New(secretKey, expiration)

	is := is.New(t)
	claims, err := authApp.ValidateJWT(encodedToken)
	is.NoErr(err)

	is.Equal("UserName", claims.Username)
	is.Equal("AuthService", claims.Issuer)
}
