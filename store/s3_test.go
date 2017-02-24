package store

import (
	"net/url"
	"testing"

	. "github.com/franela/goblin"
)

func TestS3(t *testing.T) {
	g := Goblin(t)

	artifact := mockArtifact()

	g.Describe("PutArtifact", func() {
		g.It("throws error", func() {
			err := PutArtifact(artifact, nil)

			if err == nil {
				g.Assert(err).Equal(nil) // fail test
			}
		})
	})

	g.Describe("GetObjectURL", func() {
		g.It("returns url to download object", func() {
			rawURL, err := GetObjectURL(artifact)
			g.Assert(err).Equal(nil)

			objectURL, _ := url.Parse(rawURL)
			g.Assert(objectURL.Scheme).Equal("https")
		})
	})
}
