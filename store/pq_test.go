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
		buildID = "foo"
		path    = "bar"
		key     = fakeKeyValue
	)

	return &model.Artifact{
		BuildID:   &buildID,
		Path:      &path,
		ObjectKey: &key,
	}
}

func TestDB(t *testing.T) {
	g := Goblin(t)

	db := openTestDB()
	defer db.Close()

	g.Describe("stores meta info", func() {
		g.It("saves artifact meta info", func() {
			err := db.CreateArtifact(mockArtifact())

			g.Assert(err).Equal(nil)
		})

		g.It("lists paths of artifacts", func() {
			_, err := db.ListArtifacts("foo")

			g.Assert(err).Equal(nil)
		})

		g.It("retrieves object key of an artifact", func() {
			rows, err := db.ListArtifacts("foo")

			id := rows[len(rows)-1].ID

			objectKey, err := db.RetrieveKeyOfArtifact(id, "foo")

			g.Assert(objectKey).Equal(fakeKeyValue)

			g.Assert(err).Equal(nil)
		})
	})
}
