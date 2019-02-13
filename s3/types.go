package s3

import (
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"

	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	defaultRegion = "us-east-1"
	defaultBucket = "ipfs-datastore"
	// listMax is the largest amount of objects you can request from S3 in a list
	// call.
	listMax = 1000

	// deleteMax is the largest amount of objects you can delete from S3 in a
	// delete objects call.
	deleteMax = 1000

	// used to represetn the number of concurrent batch jobs
	defaultWorkers = 100
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
	Workers       int
}

// dBatch is used to handle batch based operations
type dBatch struct {
	d       *Datastore
	ops     map[string]dBatchOp
	workers int
}

// dBatchOp is a single batch operation
type dBatchOp struct {
	val    []byte
	delete bool
}

var _ ds.Batching = (*Datastore)(nil)
