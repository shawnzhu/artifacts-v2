package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"

	"testing"

	"github.com/travis-ci/artifacts-v2/store"
	"gopkg.in/gin-gonic/gin.v1"

	. "github.com/franela/goblin"
)

func createTestApp() *gin.Engine {
	var router = gin.Default()

	router.Use(store.Store())

	router.GET("/status", HealthCheck)

	router.POST("/upload/:build_id", UploadArtifact)

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

	g.Describe("upload artifact", func() {
		g.It("responses client side error", func() {
			req, _ := http.NewRequest("POST", "/upload/bar", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusBadRequest)
		})

		g.It("fails AWS config", func() {
			var (
				bodyBuffer = bytes.Buffer{}
				bodyWriter = multipart.NewWriter(&bodyBuffer)
			)

			// file
			f, _ := ioutil.TempFile(os.TempDir(), "test")
			defer f.Close()
			f.Write([]byte("hello\nworld\n"))

			writer, _ := bodyWriter.CreateFormFile("file", f.Name())

			fh, _ := os.Open(f.Name())
			defer fh.Close()
			io.Copy(writer, fh)

			bodyWriter.Close()

			req, _ := http.NewRequest("POST", "/upload/bar", &bodyBuffer)
			req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusInternalServerError)
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

		g.It("provides download URL of artifact", func() {
			var data map[string]string

			artifactPath := fmt.Sprintf("/b/foo/a/%v", artifactID)

			req, _ := http.NewRequest("GET", artifactPath, nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusOK)

			err := json.Unmarshal(resp.Body.Bytes(), &data)
			g.Assert(err).Equal(nil)

			objectURL, _ := url.Parse(data["location"])
			g.Assert(objectURL.Scheme).Equal("https")
		})

		g.It("responses client side error", func() {
			req, _ := http.NewRequest("GET", "/b/foo/a/not-int", nil)
			resp := httptest.NewRecorder()
			app.ServeHTTP(resp, req)

			g.Assert(resp.Code).Equal(http.StatusBadRequest)
		})
	})
}
