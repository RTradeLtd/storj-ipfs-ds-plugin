package gateway

import (
	minio "github.com/minio/minio-go"
)

const (
	defaultBucket = "us-east-1"
)

// Datastore is our interface to minio
type Datastore struct {
	Client *minio.Client
	Bucket string
}
