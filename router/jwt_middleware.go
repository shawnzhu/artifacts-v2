package router

// this is a temp fix for https://github.com/gin-gonic/contrib/blob/master/jwt/jwt.go with jwt-go v3 API

import (
	"net/http"

	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"

	"github.com/gin-gonic/gin"
)

// JWT is a middleware authentiates requests via JWT
func JWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
			b := ([]byte(secret))
			return b, nil
		})

		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
	}
}
