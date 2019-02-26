package s3

import (
	"bytes"
	"fmt"
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"
	dsq "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore/query"
	"io/ioutil"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/gommon/log"
)

// NewDatastore is used to create our datastore against the minio gateway powered by storj
func NewDatastore(cfg Config) (*Datastore, error) {
	log.Debug("using config", cfg)
	// Configure to use Minio Server
	s3Config := &aws.Config{
		// TODO: determine if we need session token
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:         aws.String(cfg.Endpoint),
		Region:           aws.String(cfg.Region),
		DisableSSL:       aws.Bool(cfg.Secure),
		S3ForcePathStyle: aws.Bool(true),
	}
	s3Session, err := session.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	d := &Datastore{
		Config: cfg,
		S3:     s3.New(s3Session),
	}
	return d, nil
}

// IPFS DATASTORE FUNCTION CALLS

// Put is used to store some data
func (d *Datastore) Put(k ds.Key, value []byte) error {
	log.Info("putting object")
	resp, err := d.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
		Body:   bytes.NewReader(value),
	})
	if err != nil {
		log.Error("failed to put object", err)
		return parseError(err)
	}
	log.Info("successfully put object")
	log.Debug(resp.GoString())
	return nil
}

// Get is used to retrieve data from our storj backed s3 datastore
func (d *Datastore) Get(k ds.Key) ([]byte, error) {
	log.Info("getting object")
	resp, err := d.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
	})
	if err != nil {
		log.Error("failed to get object", err)
		return nil, parseError(err)
	}
	log.Info("successfully got object")
	log.Debug(resp.GoString())
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Has is used to check if we already have an object matching this key
func (d *Datastore) Has(k ds.Key) (exists bool, err error) {
	log.Info("checking if object exists in datastore")
	_, err = d.GetSize(k)
	if err != nil {
		log.Error("failed to check if object exists", err)
		if err == ds.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	log.Info("object exists")
	return true, nil
}

// GetSize is used to retrieve the size of an object
func (d *Datastore) GetSize(k ds.Key) (size int, err error) {
	log.Info("getting object size")
	resp, err := d.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
	})
	if err != nil {
		log.Error("failed to get object size", err)
		if s3Err, ok := err.(awserr.Error); ok && s3Err.Code() == "NotFound" {
			return -1, ds.ErrNotFound
		}
		return -1, err
	}
	log.Info("successfully got object size")
	log.Debug(resp.GoString())
	return int(*resp.ContentLength), nil
}

// Delete is used to remove an object from our datastore
func (d *Datastore) Delete(k ds.Key) error {
	log.Info("deleting object")
	resp, err := d.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
	})
	if err != nil {
		log.Error("failed to delete object", err)
		return parseError(err)
	}
	log.Info("successfully deleted object")
	log.Debug(resp.GoString())
	return nil
}

// Query is used to examine our s3 datastore and pull any objects
// matching our given query
func (d *Datastore) Query(q dsq.Query) (dsq.Results, error) {
	log.Info("executing query")
	if q.Orders != nil || q.Filters != nil {
		return nil, fmt.Errorf("storj: filters or orders are not supported")
	}

	limit := q.Limit + q.Offset
	// disabling this makes tests fail, so we should
	// investigate what exactly disabling this does
	if limit == 0 || limit > listMax {
		limit = listMax
	}
	resp, err := d.S3.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket:  aws.String(d.Bucket),
		Prefix:  aws.String(d.s3Path(q.Prefix)),
		MaxKeys: aws.Int64(int64(limit)),
	})
	if err != nil {
		log.Error("failed to list objects", err)
		return nil, err
	}

	index := q.Offset
	nextValue := func() (dsq.Result, bool) {
		for index >= len(resp.Contents) {
			if !*resp.IsTruncated {
				return dsq.Result{}, false
			}

			index -= len(resp.Contents)

			resp, err = d.S3.ListObjectsV2(&s3.ListObjectsV2Input{
				Bucket:            aws.String(d.Bucket),
				Prefix:            aws.String(d.s3Path(q.Prefix)),
				Delimiter:         aws.String("/"),
				MaxKeys:           aws.Int64(listMax),
				ContinuationToken: resp.NextContinuationToken,
			})
			if err != nil {
				log.Error("failed to list objects", err)
				return dsq.Result{Error: err}, false
			}
		}

		entry := dsq.Entry{
			Key: ds.NewKey(*resp.Contents[index].Key).String(),
		}
		if !q.KeysOnly {
			value, err := d.Get(ds.NewKey(entry.Key))
			if err != nil {
				log.Error("failed to get objects", err)
				return dsq.Result{Error: err}, false
			}
			entry.Value = value
		}

		index++
		return dsq.Result{Entry: entry}, true
	}

	return dsq.ResultsFromIterator(q, dsq.Iterator{
		Close: func() error {
			return nil
		},
		Next: nextValue,
	}), nil
}

// Close is needed to satisfy the datastore interface
func (d *Datastore) Close() error {
	return nil
}

// Batch is a batched datastore operations
func (d *Datastore) Batch() (ds.Batch, error) {
	log.Info("returning batch operation handler")
	return &dBatch{
		d:       d,
		ops:     make(map[string]dBatchOp),
		workers: d.Workers,
	}, nil
}

// TODO: not sure if we need this, borrowing from the go-s3-ds ipfs repo
func (d *Datastore) s3Path(p string) string {
	return path.Join(d.RootDirectory, p)
}

// bubble up the error, otherwise it will return nil
func parseError(err error) error {
	if s3Err, ok := err.(awserr.Error); ok && s3Err.Code() == s3.ErrCodeNoSuchKey {
		return ds.ErrNotFound
	}
	return err
}
