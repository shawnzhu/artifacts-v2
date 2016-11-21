package server

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/travis-ci/artifacts-v2/model"
	"github.com/travis-ci/artifacts-v2/store"
)

// HealthCheck provides runtime status
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

// UploadArtifact uploads an artifact file
func UploadArtifact(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	defer file.Close()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	buildID := c.Param("build_id")
	out, err := ioutil.TempFile(os.TempDir(), buildID)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	defer out.Close()
	_, err = io.Copy(out, file)

	filename := header.Filename
	path := c.PostForm("path")
	objectKey := store.HashKey(buildID, path)

	artifact := &model.Artifact{
		BuildID:   &buildID,
		Path:      &path,
		ObjectKey: &objectKey,
	}

	err = store.PutArtifact(artifact, out.Name())

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		datastore := store.FromContext(c)

		err = datastore.CreateArtifact(artifact)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			c.String(http.StatusCreated, filename)
		}
	}
}

// ListArtifacts lists artifact meta info
func ListArtifacts(c *gin.Context) {
	buildID := c.Param("build_id")

	datastore := store.FromContext(c)

	list, err := datastore.ListArtifacts(buildID)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, list)
}

// GetArtifact redirects request to a pre-signed URL of artifact file
func GetArtifact(c *gin.Context) {
	var (
		buildID    = c.Param("build_id")
		artifactID = c.Param("artifact_id")
	)

	if id, err := strconv.Atoi(artifactID); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	} else {
		datastore := store.FromContext(c)

		objectKey, err := datastore.RetrieveKeyOfArtifact(id, buildID)

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			rawURL, _ := store.GetObjectURL(objectKey)

			c.Redirect(http.StatusFound, rawURL)
		}
	}
}
