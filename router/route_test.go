package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
)

func createTestApp() http.Handler {
	return Routes()
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
	})
}
