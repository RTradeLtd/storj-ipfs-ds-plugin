package gateway

import (
	"os"
	"testing"
)

var (
	// default minio configs from our testenv
	accessKey = "C03T49S17RP0APEZDK6M"
	secretKey = "q4I9t2MN/6bAgLkbF6uyS7jtQrXuNARcyrm2vvNA"
	// these will need to be populated with the entries
	// generated from spinning up the storj sim network
	storjAccessKey = os.Getenv("STORJ_ACCESS_KEY")
	storjSecretKey = os.Getenv("STORJ_SECRET_KEY")
)

func Test_New_Config(t *testing.T) {

	sConfig := NewConfig(accessKey, secretKey)
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

func Test_New_Datastore(t *testing.T) {
	if !testing.Short() {
		accessKey = storjAccessKey
		secretKey = storjSecretKey
	}
	cfg := NewConfig(accessKey, secretKey)
	d, err := NewDatastore(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.DeleteBucket(d.Config.Bucket); err != nil {
		t.Fatal(err)
	}
	if err := d.CreateBucket(d.Config.Bucket); err != nil {
		t.Fatal(err)
	}
	if err := d.DeleteBucket(d.Config.Bucket); err != nil {
		t.Fatal(err)
	}
}
