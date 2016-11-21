package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/travis-ci/artifacts-v2/store"

	. "github.com/franela/goblin"
)

func createTestApp() *gin.Engine {
	var router = gin.Default()

	router.Use(store.Store())

	router.GET("/status", HealthCheck)

	// router.POST("/upload/:build_id", UploadArtifact)

	router.GET("/b/:build_id", ListArtifacts)

	router.GET("/b/:build_id/a/:artifact_id", GetArtifact)

	return router
}

func TestHandlers(t *testing.T) {
	g := Goblin(t)
	app := createTestApp()

	g.Describe("GET health check", func() {
		g.It("respond OK", func() {
			var data map[string]string

			req, _ := http.NewRequest("GET", "/status", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)

			err := json.Unmarshal(resp.Body.Bytes(), &data)

			g.Assert(err).Equal(nil)
			g.Assert(data["message"]).Equal("OK")
		})
	})

	g.Describe("GET meta info of artifacts", func() {
		var artifactID interface{}

		g.It("lists artifacts", func() {
			var artifacts []map[string]interface{}

			req, _ := http.NewRequest("GET", "/b/foo", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)

			err := json.Unmarshal(resp.Body.Bytes(), &artifacts)
			g.Assert(err).Equal(nil)

			artifactID = artifacts[0]["ID"]
		})

		g.It("redirects to download artifact", func() {
			artifactPath := fmt.Sprintf("/b/foo/a/%v", artifactID)
			fmt.Println(artifactPath)
			req, _ := http.NewRequest("GET", artifactPath, nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusFound)

			rawURL := resp.HeaderMap.Get("Location")
			objectURL, _ := url.Parse(rawURL)
			g.Assert(objectURL.Scheme).Equal("https")
		})
	})
}
