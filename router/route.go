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

	buildBase := mux.NewRouter()

	router.PathPrefix("/builds/{build_id}").Handler(n.With(
		negroni.Wrap(buildBase),
	))

	buildBase.Methods("GET").HandlerFunc(server.ListArtifacts)
	buildBase.Methods("POST").HandlerFunc(server.UploadArtifact)
	buildBase.Methods("GET").Path("/artifacts/{artifact_id}").HandlerFunc(server.GetArtifact)

	return router
}
