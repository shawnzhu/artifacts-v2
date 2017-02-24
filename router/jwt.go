package router

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/auth0/go-jwt-middleware"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

// JWT returns a new negroni middleware instance authenticates request
// validating JWT signature by a given public key.
func JWT(c *cli.Context) negroni.HandlerFunc {

	publicKeyPEM := []byte(c.String("jwt-public-key"))

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
	}

	m := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: keyFunc,
		SigningMethod:       jwt.SigningMethodRS256,
	})

	return negroni.HandlerFunc(m.HandlerWithNext)
}
