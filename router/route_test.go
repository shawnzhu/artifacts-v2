package router

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/urfave/cli"

	. "github.com/franela/goblin"
)

func createTestApp() http.Handler {
	mockSet := flag.NewFlagSet("test", 0)
	mockSet.String("jwt-public-key", os.Getenv("JWT_PUBLIC_KEY"), "")
	return Routes(cli.NewContext(nil, mockSet, nil))
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

	g.Describe("URI /status", func() {
		g.It("no token required", func() {
			req, _ := http.NewRequest("GET", "/status", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)
		})
	})

	g.Describe("URI /jobs/foo", func() {

		g.It("supports verb OPTIONS", func() {
			req, _ := http.NewRequest("OPTIONS", "/jobs/foo", nil)
			req.Header.Set("Origin", "http://foo.example.com")
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)
		})

		g.It("requires token", func() {
			req, _ := http.NewRequest("GET", "/jobs/foo", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusUnauthorized)
		})

		g.It("supports JWT token", func() {
			jwtToken, _ := generateJWTToken()

			req, _ := http.NewRequest("GET", "/jobs/foo", nil)
			req.Header.Set("Authorization", "BEARER "+jwtToken)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)
		})

		g.It("supports CORS headers", func() {
			jwtToken, _ := generateJWTToken()

			req, _ := http.NewRequest("GET", "/jobs/foo", nil)
			req.Header.Set("Authorization", "BEARER "+jwtToken)
			req.Header.Set("Origin", "http://foo.example.com")
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)
			g.Assert(resp.Header().Get("Access-Control-Allow-Credentials")).Equal("true")
		})
	})
}
