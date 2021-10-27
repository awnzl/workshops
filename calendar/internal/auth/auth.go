package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type App struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JWTClaim struct {
	Username string
	jwt.StandardClaims
}

func New(k string, h int64) *App {
	return &App{
		SecretKey:       k,
		Issuer:          "AuthService",
		ExpirationHours: h,
	}
}

func (a *App) GenerateJWT(username string) (signedToken string, err error) {
	claims := &JWTClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(a.ExpirationHours)).Unix(),
			Issuer:    a.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(a.SecretKey))

	return
}

func (a *App) ValidateJWT(signedToken string) (claims *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return
}
