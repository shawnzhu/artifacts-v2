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

	jobID := mux.Vars(r)["job_id"]
	filename := header.Filename

	artifact := &model.Artifact{
		JobID: &jobID,
		Path:  &filename,
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
	jobID := mux.Vars(r)["job_id"]

	datastore := store.FromContext(r)

	list, err := datastore.ListArtifacts(jobID)

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

	var artifactID = vars["artifact_id"]

	if id, err := strconv.Atoi(artifactID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		datastore := store.FromContext(r)

		artifact, err := datastore.RetrieveArtifact(id)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			rawURL, _ := store.GetObjectURL(artifact)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf("{\"location\": \"%s\"}", rawURL)))
		}
	}
}
