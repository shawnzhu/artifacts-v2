package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/travis-ci/artifacts-v2/model"
	"github.com/travis-ci/artifacts-v2/store"
)

// HealthCheck provides runtime status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"message\":\"OK\"}"))
}

// UploadArtifact uploads an artifact file
func UploadArtifact(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if file != nil {
		defer file.Close()
	}

	buildID := mux.Vars(r)["job_id"]

	filename := header.Filename
	objectKey := store.HashKey(buildID, filename)

	artifact := &model.Artifact{
		BuildID:   &buildID,
		Path:      &filename,
		ObjectKey: &objectKey,
	}

	err = store.PutArtifact(artifact, file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datastore := store.FromContext(r)

	err = datastore.CreateArtifact(artifact)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(filename))
	}
}

// ListArtifacts lists artifact meta info
func ListArtifacts(w http.ResponseWriter, r *http.Request) {
	buildID := mux.Vars(r)["job_id"]

	datastore := store.FromContext(r)

	list, err := datastore.ListArtifacts(buildID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		data, _ := json.Marshal(list)
		w.Write(data)
	}
}

// GetArtifact redirects request to a pre-signed URL of artifact file
func GetArtifact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var (
		buildID    = vars["job_id"]
		artifactID = vars["artifact_id"]
	)

	if id, err := strconv.Atoi(artifactID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		datastore := store.FromContext(r)

		objectKey, err := datastore.RetrieveKeyOfArtifact(id, buildID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			rawURL, _ := store.GetObjectURL(objectKey)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"location\": \"%s\"}", rawURL)))
		}
	}
}
