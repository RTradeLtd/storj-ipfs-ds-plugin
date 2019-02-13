package s3

import "github.com/aws/aws-sdk-go/service/s3"

const (
	defaultRegion = "us-east-1"
	defaultBucket = "ipfs-datastore"
	// listMax is the largest amount of objects you can request from S3 in a list
	// call.
	listMax = 1000

	// deleteMax is the largest amount of objects you can delete from S3 in a
	// delete objects call.
	deleteMax = 1000
)

// Datastore is our interface to minio
type Datastore struct {
	S3 *s3.S3
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

// NewConfig is used to generate a config with defaults
func NewConfig(accessKey, secretKey string) Config {
	return Config{
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Bucket:        defaultBucket,
		Region:        defaultRegion,
		Endpoint:      "http://127.0.0.1:9000",
		RootDirectory: "",
		Secure:        false,
	}
}

type dBatch struct {
	d   *Datastore
	ops map[string]dBatchOp
}

type dBatchOp struct {
	val    []byte
	delete bool
}
