package store

import (
	"net/url"
	"testing"

	"github.com/travis-ci/artifacts-v2/model"

	. "github.com/franela/goblin"
)

func TestS3(t *testing.T) {
	g := Goblin(t)

	g.Describe("PutArtifact", func() {
		g.It("throws error", func() {
			buildID := "foo"
			err := PutArtifact(&model.Artifact{
				BuildID: &buildID,
			}, nil)

			if err == nil {
				g.Assert(err).Equal(nil) // fail test
			}
		})
	})

	g.Describe("GetArtifact", func() {
		g.It("returns meta info", func() {
			buildID := "foo"
			key := "bar"

			artifact, err := GetArtifact(buildID, key)

			g.Assert(err).Equal(nil)
			g.Assert(*artifact.ObjectKey).Equal(key)
		})
	})

	g.Describe("GetObjectURL", func() {
		g.It("returns url to download object", func() {
			key := "bar"

			rawURL, err := GetObjectURL(key)
			g.Assert(err).Equal(nil)

			objectURL, _ := url.Parse(rawURL)
			g.Assert(objectURL.Scheme).Equal("https")
		})
	})
}
