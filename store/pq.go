package store

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/urfave/cli"
	"github.com/urfave/negroni"

	"github.com/travis-ci/artifacts-v2/model"

	// db driver
	_ "github.com/lib/pq"
)

type datastore struct {
	*sql.DB
}

// The key type is unexported to prevent collisions with context keys defined in
// other packages
type key int

const (
	defaultDBURL    string = "postgresql://postgres@localhost:5432/test_artifacts?sslmode=disable"
	storeContextKey key    = 0
)

// open opens new db connection
func open(driverName, dbConnURL string) *datastore {
	db, err := sql.Open(driverName, dbConnURL)

	if err != nil {
		log.Fatal(err)
	}

	return &datastore{db}
}

// WithStore mixins datastore into request context
func WithStore(c *cli.Context) negroni.HandlerFunc {
	var (
		dbURL string
		store *datastore
	)

	if dbURL = c.String("db-url"); dbURL == "" {
		dbURL = defaultDBURL
	}

	store = open("postgres", dbURL)

	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := context.WithValue(r.Context(), storeContextKey, store)
		next.ServeHTTP(rw, r.WithContext(ctx))
	}
}

// FromContext extract datastore from given context
func FromContext(r *http.Request) *datastore {
	if store := r.Context().Value(storeContextKey); store != nil {
		return store.(*datastore)
	}

	return nil
}

// CreateArtifact is for saving meta info
func (db *datastore) CreateArtifact(artifact *model.Artifact) error {
	_, err := db.Exec(`INSERT INTO artifacts_v2.artifacts (job_id, path)
		VALUES ($1, $2)`, artifact.JobID, artifact.Path)

	return err
}

func (db *datastore) RetrieveArtifact(artifactID int) (*model.Artifact, error) {
	var (
		jobID string
		path  string
	)

	err := db.QueryRow(`SELECT job_id, path FROM artifacts_v2.artifacts
		WHERE artifact_id = $1`, artifactID).Scan(&jobID, &path)

	return &model.Artifact{
		ID:    artifactID,
		JobID: &jobID,
		Path:  &path,
	}, err
}

func (db *datastore) ListArtifacts(jobID string) ([]*model.Artifact, error) {
	rows, err := db.Query(`SELECT artifact_id, path FROM artifacts_v2.artifacts
		WHERE job_id = $1`, jobID)

	if err != nil {
		log.Fatal(err)
	}

	artifacts := []*model.Artifact{}

	defer rows.Close()
	for rows.Next() {
		var (
			id   int
			path string
		)

		err = rows.Scan(&id, &path)

		if err == nil {
			artifacts = append(artifacts, &model.Artifact{
				ID:    id,
				JobID: &jobID,
				Path:  &path,
			})
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return artifacts, err
}
