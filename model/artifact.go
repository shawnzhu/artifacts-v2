package model

// Artifact represents metainfo of user artifact
type Artifact struct {
	ID      int     `db:"artifact_id,pk"`
	BuildID *string `db:"build_id"` // deprecated
	JobID   *string `db:"job_id"`
	Path    *string `db:"path"`
}
