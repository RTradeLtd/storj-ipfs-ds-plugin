package s3

import (
	dstest "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore/test"
	"os"
	"testing"
)

var (
	// default minio configs from our testenv
	accessKey = "C03T49S17RP0APEZDK6M"
	secretKey = "q4I9t2MN/6bAgLkbF6uyS7jtQrXuNARcyrm2vvNA"
	logPath   = "./tmp"
	// these will need to be populated with the entries
	// generated from spinning up the storj sim network
	storjAccessKey = os.Getenv("STORJ_ACCESS_KEY")
	storjSecretKey = os.Getenv("STORJ_SECRET_KEY")
)

func Test_New_Config(t *testing.T) {

	sConfig := NewConfig(accessKey, secretKey, logPath)
	if sConfig.AccessKey != accessKey {
		t.Fatal("failed to set correct access key")
	}
	if sConfig.SecretKey != secretKey {
		t.Fatal("failed to set correct secret key")
	}
	if sConfig.Bucket != defaultBucket {
		t.Fatal("failed to set correct bucket")
	}
	if sConfig.Region != defaultRegion {
		t.Fatal("failed to set correct region")
	}
}

func Test_Datastore_Non_Storj(t *testing.T) {
	cfg := NewConfig(accessKey, secretKey, logPath)
	d, err := NewDatastore(cfg, true)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.CreateBucket(d.Config.Bucket); err != nil {
		t.Fatal(err)
	}
	if err := d.BucketExists(d.Config.Bucket); err != nil {
		t.Fatal(err)
	}
	if err := d.BucketExists("randombucketname"); err == nil {
		t.Fatal("expected error")
	}
	defer d.DeleteBucket(d.Config.Bucket)

	t.Run("basic operations", func(t *testing.T) {
		dstest.SubtestBasicPutGet(t, d)
	})

	t.Run("not found operations", func(t *testing.T) {
		dstest.SubtestNotFounds(t, d)
	})

	t.Run("many puts and gets, query", func(t *testing.T) {
		dstest.SubtestManyKeysAndQuery(t, d)
	})
}
