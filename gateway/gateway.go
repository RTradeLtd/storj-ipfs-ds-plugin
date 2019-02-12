package gateway

import (
	"bytes"
	ds "gx/ipfs/QmUadX5EcvrBmxAV9sE7wUWtWSqxns5K84qKJBixmcT1w9/go-datastore"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// NewDatastore is used to create our datastore against the minio gateway powered by storj
func NewDatastore(cfg Config) (*Datastore, error) {
	// Configure to use Minio Server
	s3Config := &aws.Config{
		// TODO: determine if we need session token
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		DisableSSL:       aws.Bool(cfg.Secure),
		S3ForcePathStyle: aws.Bool(true),
	}
	s3Session := session.New(s3Config)
	s3Client := s3.New(s3Session)
	createParam := &s3.CreateBucketInput{
		Bucket: aws.String(cfg.Bucket),
	}

	if _, err := s3Client.CreateBucket(createParam); err != nil {
		return nil, err
	}
	return &Datastore{
		Config: cfg,
		Store:  s3Client,
	}, nil
}

// NewConfig is used to generate a config with defaults
func NewConfig(accessKey, secretKey string) Config {
	return Config{
		AccessKey:     accessKey,
		SecretKey:     secretKey,
		Bucket:        defaultBucket,
		Endpoint:      "http://127.0.0.1:9000",
		RootDirectory: "",
		Secure:        false,
	}
}

// Put is used to store some data
func (d *Datastore) Put(k ds.Key, value []byte) error {
	_, err := d.Store.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
		Body:   bytes.NewReader(value),
	})
	return err
}

// TODO: not sure if we need this, borrowing from the go-s3-ds ipfs repo
func (d *Datastore) s3Path(p string) string {
	return path.Join(d.RootDirectory, p)
}
