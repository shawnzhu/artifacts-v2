package router

import (
	"os"

	"github.com/dgrijalva/jwt-go"

	"github.com/auth0/go-jwt-middleware"
	"github.com/urfave/negroni"
)

// JWT returns a new negroni middleware instance authenticates request
// validating JWT signature by a given public key.
func JWT() negroni.HandlerFunc {
	publicKeyPEM := []byte(os.Getenv("JWT_PUBLIC_KEY"))

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	}

	m := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: keyFunc,
		SigningMethod:       jwt.SigningMethodRS256,
	})

	return negroni.HandlerFunc(m.HandlerWithNext)
}
