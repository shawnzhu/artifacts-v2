package store

import (
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"crypto/md5"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/travis-ci/artifacts-v2/model"
)

const (
	artifactBucketName = "travis-ci-artifacts-test"
	downloadExpireTime = 5 * time.Minute // duration of download URL
)

func getBucketName() *string {
	var bucketName = os.Getenv("ARTIFACTS_S3_BUCKET_NAME")

	if bucketName == "" {
		bucketName = artifactBucketName
	}

	return aws.String(bucketName)
}

func newAWSSession() (*s3.S3, error) {
	// TODO make region configurable
	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})

	svc := s3.New(sess)

	return svc, err
}

// ObjectKey generates object key from artifact meta info
func ObjectKey(artifact *model.Artifact) string {
	data := []byte(fmt.Sprintf("%s-%s", artifact.JobID, artifact.Path))
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

// PutArtifact saves artifact to S3
func PutArtifact(artifact *model.Artifact, file multipart.File) error {
	contentDispositionHeader := fmt.Sprintf("attachment; filename=%s",
		aws.StringValue(artifact.Path))
	svc, err := newAWSSession()

	objectKey := ObjectKey(artifact)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:             getBucketName(),
		ContentDisposition: &contentDispositionHeader,
		Key:                &objectKey,
		Body:               file,
	})

	return err
}

// GetArtifact retrieves artifact content
func GetArtifact(jobID string, key string) (*model.Artifact, error) {

	return &model.Artifact{
		JobID: &jobID,
		Path:  &key,
	}, nil
}

// GetObjectURL returns a download URL of an S3 object
func GetObjectURL(artifact *model.Artifact) (string, error) {
	svc, _ := newAWSSession()
	objectKey := ObjectKey(artifact)

	getObjectReq, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: getBucketName(),
		Key:    aws.String(objectKey),
	})

	return getObjectReq.Presign(downloadExpireTime)
}
