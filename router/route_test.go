package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"

	. "github.com/franela/goblin"
)

func createTestApp() http.Handler {
	return Routes()
}

func generateJWTToken() (string, error) {
	rsaPrivateKey, _ := ioutil.ReadFile(os.Getenv("JWT_PRIVATE_KEY_PATH"))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(rsaPrivateKey)
	jwtClaims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtClaims)

	return token.SignedString(privateKey)
}

func TestJWT(t *testing.T) {
	g := Goblin(t)
	app := createTestApp()

	g.Describe("GET /builds/foo", func() {

		g.It("requires token", func() {
			req, _ := http.NewRequest("GET", "/builds/foo", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusUnauthorized)
		})

		g.It("supports JWT token", func() {
			jwtToken, _ := generateJWTToken()

			req, _ := http.NewRequest("GET", "/builds/foo", nil)
			req.Header.Set("Authorization", "BEARER "+jwtToken)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)
		})
	})
}
