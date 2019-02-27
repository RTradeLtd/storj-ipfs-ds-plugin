package storj

import (
	"testing"

	"github.com/RTradeLtd/storj-ipfs-ds-plugin/s3"
)

const (
	expectedType          = "storj"
	expectedName          = "ds-storj"
	expectedVersion       = "v0.0.0"
	expectedDatastoreType = "storj"
	accessKey             = "accessKey123"
	secretKey             = "secretKey123"
	pluginTestBucket      = "plugin-test-bucket"
)

func Test_Plugin(t *testing.T) {
	sp := &SJPlugin{}
	if sp.Name() != expectedName {
		t.Fatal("bad name returned")
	}
	if sp.Version() != expectedVersion {
		t.Fatal("bad version returned")
	}
	if err := sp.Init(); err != nil {
		t.Fatal(err)
	}
	if sp.DatastoreTypeName() != expectedType {
		t.Fatal("bad datastore type returned")
	}
	cfg := s3.NewConfig(accessKey, secretKey, "")
	cfg.Bucket = pluginTestBucket
	dsc := DSConfig{cfg}
	if _, err := dsc.Create("mypath"); err != nil {
		t.Fatal(err)
	}
}
