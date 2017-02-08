package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/travis-ci/artifacts-v2/server"
	"github.com/travis-ci/artifacts-v2/store"
)

// Routes load middlewares
func Routes() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/status", server.HealthCheck).Methods("GET")

	n := negroni.New(
		JWT(),
		CORS(),
		store.WithStore(),
	)

	router.Methods("POST").Path("/upload/{build_id}").HandlerFunc(server.UploadArtifact)

	build := router.PathPrefix("/builds/{build_id}").Subrouter()

	build.Methods("GET").HandlerFunc(server.ListArtifacts)
	build.Methods("POST").HandlerFunc(server.UploadArtifact)
	build.Methods("GET").Path("/artifacts/{artifact_id}").HandlerFunc(server.GetArtifact)

	n.UseHandler(router)

	return n
}
