package gateway

import (
	"bytes"
	ds "gx/ipfs/QmUadX5EcvrBmxAV9sE7wUWtWSqxns5K84qKJBixmcT1w9/go-datastore"

	minio "github.com/minio/minio-go"
)

// NewDataStore is used to initialize our connection to minio, creating our datastore wrapper
func NewDataStore(endpoint, accessKeyID, secretAccessKey string, secure bool) (*Datastore, error) {
	// connect to minio
	mini, err := minio.New(endpoint, accessKeyID, secretAccessKey, secure)
	if err != nil {
		return nil, err
	}
	// verify our connection
	if _, err := mini.ListBuckets(); err != nil {
		return nil, err
	}
	// return our datastore wrapper
	return &Datastore{
		Client: mini,
		Bucket: defaultBucket,
	}, nil
}

// Put is used to store an object in our minio backend connect to storj
func (d *Datastore) Put(k ds.Key, value []byte) error {
	_, err := d.Client.PutObject(d.Bucket, k.Name(), bytes.NewReader(value), int64(len(value)), minio.PutObjectOptions{})
	return err
}

/*
func (s *S3Bucket) Put(k ds.Key, value []byte) error {
	_, err := s.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.s3Path(k.String())),
		Body:   bytes.NewReader(value),
	})
	return parseError(err)
}
*/
