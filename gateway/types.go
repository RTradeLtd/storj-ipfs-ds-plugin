package gateway

import "github.com/aws/aws-sdk-go/service/s3"

const (
	defaultRegion = "us-east-1"
	defaultBucket = "ipfs-datastore"
)

// Datastore is our interface to minio
type Datastore struct {
	Store *s3.S3
	Config
}

// Config is used to configure our gateway
type Config struct {
	AccessKey string
	SecretKey string
	//	SessionToken   string
	Bucket        string
	Region        string
	Endpoint      string
	RootDirectory string
	Secure        bool
}
