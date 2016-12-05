package router

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

// Auth authenticates request validating JWT signature by a given public key.
// a user agent must provide a BEARER token via the Authorization header
func Auth() gin.HandlerFunc {
	publicKeyPEM, _ := ioutil.ReadFile(os.Getenv("JWT_PUBLIC_KEY_PATH"))

	return func(c *gin.Context) {
		keyFunc := func(token *jwt.Token) (interface{}, error) {
			return jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
		}

		_, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, keyFunc)

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
	}
}
