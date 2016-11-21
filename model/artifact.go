package model

// Artifact represents metainfo of user artifact
type Artifact struct {
	ID        int     `db:"artifact_id,pk"`
	BuildID   *string `db:"build_id"`
	Path      *string `db:"path"`
	ObjectKey *string `db:"s3_object_key"`
}
