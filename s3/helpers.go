package s3

import (
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func worker(jobs <-chan func() error, results chan<- error) {
	for j := range jobs {
		results <- j()
	}
}

// BucketExists is used to lookup if the designated bucket exists
func (d *Datastore) BucketExists(name string) error {
	listParam := &s3.ListBucketsInput{}
	out, err := d.S3.ListBuckets(listParam)
	if err != nil {
		return parseError(err)
	}
	for _, v := range out.Buckets {
		if *v.Name == name {
			return nil
		}
	}
	return ds.ErrNotFound
}

// CreateBucket is used to create a bucket
func (d *Datastore) CreateBucket(name string) error {
	createParam := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	// create bucket ensure we have initialize client correct
	_, err := d.S3.CreateBucket(createParam)
	return parseError(err)
}

// DeleteBucket is used to remove the specified bucket
func (d *Datastore) DeleteBucket(name string) error {
	deleteParam := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err := d.S3.DeleteBucket(deleteParam)
	return parseError(err)
}

func (d *Datastore) logDebug(err error) {
	if d.debugLogging == true {
		d.l.Error(err)
	}
}
