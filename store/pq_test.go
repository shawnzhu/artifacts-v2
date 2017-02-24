package store

import (
	"testing"

	"github.com/travis-ci/artifacts-v2/model"

	. "github.com/franela/goblin"
)

const fakeKeyValue = "fake_key"

func openTestDB() *datastore {
	return open("postgres", "postgresql://postgres@localhost:5432/test_artifacts?sslmode=disable")
}

func mockArtifact() *model.Artifact {
	var (
		jobID = "foo"
		path  = "bar"
	)

	return &model.Artifact{
		JobID: &jobID,
		Path:  &path,
	}
}

func TestDB(t *testing.T) {
	g := Goblin(t)

	db := openTestDB()
	defer db.Close()

	artifact := mockArtifact()

	g.Describe("stores meta info", func() {
		g.It("saves artifact meta info", func() {
			err := db.CreateArtifact(artifact)

			g.Assert(err).Equal(nil)
		})

		g.It("lists paths of artifacts", func() {
			_, err := db.ListArtifacts("foo")

			g.Assert(err).Equal(nil)
		})

		g.It("retrieves an artifact", func() {
			rows, err := db.ListArtifacts(*artifact.JobID)

			id := rows[len(rows)-1].ID

			firstArtifact, _ := db.RetrieveArtifact(id)

			g.Assert(firstArtifact.JobID).Equal(artifact.JobID)

			g.Assert(err).Equal(nil)
		})
	})
}
