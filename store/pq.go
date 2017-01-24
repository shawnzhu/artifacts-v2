package store

import (
	"database/sql"
	"log"
	"os"

	"gopkg.in/gin-gonic/gin.v1"

	"github.com/travis-ci/artifacts-v2/model"

	// db driver
	_ "github.com/lib/pq"
)

type datastore struct {
	*sql.DB
}

const defaultDBURL = "postgresql://postgres@localhost:5432/test_artifacts?sslmode=disable"

// open opens new db connection
func open(driverName, dbConnURL string) *datastore {
	db, err := sql.Open(driverName, dbConnURL)

	if err != nil {
		log.Fatal(err)
	}

	return &datastore{db}
}

// Store provides a middleware to inject data source
func Store() gin.HandlerFunc {
	var dbURL string

	if dbURL = os.Getenv("DB_URL"); dbURL == "" {
		dbURL = defaultDBURL
	}

	return func(c *gin.Context) {
		var store *datastore
		store = open("postgres", dbURL)
		c.Set("store", store)
		c.Next()
	}
}

// FromContext extract datastore from given context
func FromContext(c *gin.Context) *datastore {
	return c.MustGet("store").(*datastore)
}

// CreateArtifact is for saving meta info
func (db *datastore) CreateArtifact(artifact *model.Artifact) error {
	_, err := db.Exec(`INSERT INTO artifacts (build_id, path, s3_object_key)
		VALUES ($1, $2, $3)`, artifact.BuildID, artifact.Path, artifact.ObjectKey)

	return err
}

func (db *datastore) RetrieveKeyOfArtifact(id int, buildID string) (string, error) {
	var objectKey string

	err := db.QueryRow(`SELECT s3_object_key FROM artifacts
		WHERE build_id = $1 AND artifact_id = $2`, buildID, id).Scan(&objectKey)

	return objectKey, err
}

func (db *datastore) ListArtifacts(buildID string) ([]*model.Artifact, error) {
	rows, err := db.Query(`SELECT artifact_id, path FROM artifacts
		WHERE build_id = $1`, buildID)

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
				ID:      id,
				BuildID: &buildID,
				Path:    &path,
			})
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return artifacts, err
}
