package s3

import (
	"fmt"
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/gommon/log"
)

// Put is a batch based put operation
func (db *dBatch) Put(k ds.Key, val []byte) error {
	db.ops[k.String()] = dBatchOp{
		val:    val,
		delete: false,
	}
	return nil
}

// Delete is a batch based delete operation
func (db *dBatch) Delete(k ds.Key) error {
	db.ops[k.String()] = dBatchOp{
		val:    nil,
		delete: true,
	}
	return nil
}

// Commit is used to commit batch operations and finalize their actions
func (db *dBatch) Commit() error {
	var (
		deleteObjs []*s3.ObjectIdentifier
		putKeys    []ds.Key
	)
	for k, op := range db.ops {
		if op.delete {
			deleteObjs = append(deleteObjs, &s3.ObjectIdentifier{
				Key: aws.String(k),
			})
		} else {
			putKeys = append(putKeys, ds.NewKey(k))
		}
	}

	numJobs := len(putKeys) + (len(deleteObjs) / deleteMax)
	jobs := make(chan func() error, numJobs)
	results := make(chan error, numJobs)

	numWorkers := db.workers
	if numJobs < numWorkers {
		numWorkers = numJobs
	}

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	defer wg.Wait()

	for w := 0; w < numWorkers; w++ {
		go func() {
			defer wg.Done()
			worker(jobs, results)
		}()
	}

	for _, k := range putKeys {
		jobs <- db.newPutJob(k, db.ops[k.String()].val)
	}

	if len(deleteObjs) > 0 {
		for i := 0; i < len(deleteObjs); i += deleteMax {
			limit := deleteMax
			if len(deleteObjs[i:]) < limit {
				limit = len(deleteObjs[i:])
			}

			jobs <- db.newDeleteJob(deleteObjs[i : i+limit])
		}
	}
	close(jobs)

	var errs []string
	for i := 0; i < numJobs; i++ {
		err := <-results
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("storj: failed batch operation:\n%s", strings.Join(errs, "\n"))
	}

	return nil
}

func (db *dBatch) newPutJob(k ds.Key, value []byte) func() error {
	return func() error {
		return db.d.Put(k, value)
	}
}

func (db *dBatch) newDeleteJob(objs []*s3.ObjectIdentifier) func() error {
	return func() error {
		resp, err := db.d.S3.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(db.d.Bucket),
			Delete: &s3.Delete{
				Objects: objs,
			},
		})
		if err != nil {
			log.Error("failed to execute Delete Objects", err)
			return err
		}

		var errs []string
		for _, err := range resp.Errors {
			errs = append(errs, err.String())
		}

		if len(errs) > 0 {
			return fmt.Errorf("storj: failed to delete objects: %s", errs)
		}

		return nil
	}
}
