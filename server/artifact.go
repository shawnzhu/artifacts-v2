package server

import (
	"net/http"
	"strconv"

	"github.com/travis-ci/artifacts-v2/model"
	"github.com/travis-ci/artifacts-v2/store"
	"gopkg.in/gin-gonic/gin.v1"
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

	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if file != nil {
		defer file.Close()
	}

	buildID := c.Param("build_id")

	filename := header.Filename
	objectKey := store.HashKey(buildID, filename)

	artifact := &model.Artifact{
		BuildID:   &buildID,
		Path:      &filename,
		ObjectKey: &objectKey,
	}

	err = store.PutArtifact(artifact, file)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	datastore := store.FromContext(c)

	err = datastore.CreateArtifact(artifact)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.String(http.StatusCreated, filename)
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

			c.JSON(http.StatusOK, gin.H{
				"location": rawURL,
			})
		}
	}
}
