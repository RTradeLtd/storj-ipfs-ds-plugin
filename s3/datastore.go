package s3

import (
	"bytes"
	"fmt"
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"
	dsq "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore/query"
	"io/ioutil"
	"path"

	"go.uber.org/zap"

	"github.com/RTradeLtd/storj-ipfs-ds-plugin/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Datastore is our interface to minio
type Datastore struct {
	S3 *s3.S3
	l  *zap.SugaredLogger
	Config
}

// NewDatastore is used to create our datastore against the minio gateway powered by storj
func NewDatastore(cfg Config, dev bool) (*Datastore, error) {
	logger, err := log.NewLogger(cfg.LogPath, dev)
	if err != nil {
		return nil, err
	}
	logger.Infow("initialized logger", "config", cfg)
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
		l:      logger,
	}
	return d, nil
}

// IPFS DATASTORE FUNCTION CALLS

// Put is used to store some data
func (d *Datastore) Put(k ds.Key, value []byte) error {
	d.l.Infow("putting object", "key", k)
	resp, err := d.S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
		Body:   bytes.NewReader(value),
	})
	if err != nil {
		d.l.Errorw("failed to put object", "error", err)
		return parseError(err)
	}
	d.l.Info("successfully put object")
	d.l.Debug(resp.GoString())
	return nil
}

// Get is used to retrieve data from our storj backed s3 datastore
func (d *Datastore) Get(k ds.Key) ([]byte, error) {
	d.l.Infow("getting object", "key", k)
	resp, err := d.S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
	})
	if err != nil {
		d.l.Errorw("failed to get object", "error", err)
		return nil, parseError(err)
	}
	d.l.Info("successfully got object")
	d.l.Debug(resp.GoString())
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

// Has is used to check if we already have an object matching this key
func (d *Datastore) Has(k ds.Key) (exists bool, err error) {
	d.l.Infow("checking datastore for object", "key", k)
	_, err = d.GetSize(k)
	if err != nil {
		d.l.Errorw("failed to check datastore", "error", err)
		if err == ds.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	d.l.Info("object exists")
	return true, nil
}

// GetSize is used to retrieve the size of an object
func (d *Datastore) GetSize(k ds.Key) (size int, err error) {
	d.l.Infow("getting object size", "key", k)
	resp, err := d.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
	})
	if err != nil {
		d.l.Errorw("failed to get object size", "error", err)
		if s3Err, ok := err.(awserr.Error); ok && s3Err.Code() == "NotFound" {
			return -1, ds.ErrNotFound
		}
		return -1, err
	}
	d.l.Infow("successfully got object size")
	d.l.Debug(resp.GoString())
	return int(*resp.ContentLength), nil
}

// Delete is used to remove an object from our datastore
func (d *Datastore) Delete(k ds.Key) error {
	d.l.Infow("deleting object", "key", k)
	resp, err := d.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.s3Path(k.String())),
	})
	if err != nil {
		d.l.Errorw("failed to delete object", "error", err)
		return parseError(err)
	}
	d.l.Info("successfully deleted object")
	d.l.Debug(resp.GoString())
	return nil
}

// Query is used to examine our s3 datastore and pull any objects
// matching our given query
func (d *Datastore) Query(q dsq.Query) (dsq.Results, error) {
	d.l.Infow("executing query", "query", q)
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
		d.l.Errorw("failed to list objects while running query", "error", err)
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
				d.l.Errorw("failed to list objects while running query", "error", err)
				return dsq.Result{Error: err}, false
			}
		}

		entry := dsq.Entry{
			Key: ds.NewKey(*resp.Contents[index].Key).String(),
		}
		if !q.KeysOnly {
			value, err := d.Get(ds.NewKey(entry.Key))
			if err != nil {
				d.l.Errorw("failed to get objects while running query", "error", err)
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
	d.l.Info("returning batch operation handler")
	return &dBatch{
		d:       d,
		l:       d.l.With("batch"),
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
