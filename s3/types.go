package s3

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

// Config is used to configure our gateway
type Config struct {
	AccessKey string
	SecretKey string
	//	SessionToken   string
	Bucket        string
	Region        string
	Endpoint      string
	RootDirectory string
	LogPath       string
	Secure        bool
	Workers       int
}
