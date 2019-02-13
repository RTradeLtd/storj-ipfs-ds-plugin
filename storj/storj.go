package storj

import (
	"fmt"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/plugin"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo"

	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo"

	"github.com/RTradeLtd/storj-ipfs-ds-plugin/s3"
)

// SJPlugin is used to allow storj nodes
// to function as the storage backend for a go-ipfs node.
type SJPlugin struct{}

// DatastoreType is the type of datastore name
var DatastoreType = "storj"

var _ plugin.PluginDatastore = (*SJPlugin)(nil)

// plugin.Plugin mandatory interfaces

// Name returns the plugin name
func (sp *SJPlugin) Name() string {
	return "ds-storj"
}

// Version returns the plugin version
func (sp *SJPlugin) Version() string {
	return "v0.0.0"
}

// Init is used to hold initialization logic
func (sp *SJPlugin) Init() error {
	return nil
}

// plugin.DatastorePlugin mandatory interfaces

// DatastoreTypeName returns the datastore name
// and must be unique across all ipfs datastores
func (sp *SJPlugin) DatastoreTypeName() string {
	return DatastoreType
}

// DatastoreConfigParser returns a configuration parser for storj datastore
func (sp *SJPlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(m map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		accessKey := m["accessKey"]
		if accessKey == nil {
			accessKey = ""
		}

		secretKey := m["secretKey"]
		if secretKey == nil {
			secretKey = ""
		}
		workers := m["workers"]
		if workers == nil {
			workers = 1000
		}
		bucket, ok := m["bucket"].(string)
		if !ok {
			return nil, fmt.Errorf("ds-storj: unable to convert bucket to string type")
		}
		if bucket == "" {
			return nil, fmt.Errorf("ds-storj: bucket configuration is empty")
		}

		region, ok := m["region"].(string)
		if !ok {
			return nil, fmt.Errorf("ds-storj: unable to convert region to string type")
		}
		if region == "" {
			return nil, fmt.Errorf("ds-storj: region configuration is empty")
		}

		endpoint, ok := m["endpoint"].(string)
		if !ok {
			return nil, fmt.Errorf("ds-storj: unable to convert endpoint to string type")
		}
		if endpoint == "" {
			return nil, fmt.Errorf("ds-storj: endpoint configuration is empty")
		}

		rootDirectory, ok := m["rootDirectory"].(string)
		if !ok {
			return nil, fmt.Errorf("ds-storj: unable to convert rootDirectory to string type")
		}
		// permit empty string for root directory

		return &DSConfig{
			cfg: s3.Config{
				AccessKey:     accessKey.(string),
				SecretKey:     secretKey.(string),
				Bucket:        bucket,
				Region:        region,
				Endpoint:      endpoint,
				RootDirectory: rootDirectory,
				Workers:       workers.(int),
			},
		}, nil
	}
}

// DSConfig is the configuration for our datastore
type DSConfig struct {
	cfg s3.Config
}

// DiskSpec returns the disk specification of our s3 datastore
func (dsc *DSConfig) DiskSpec() fsrepo.DiskSpec {
	return map[string]interface{}{
		"bucket":        dsc.cfg.Bucket,
		"region":        dsc.cfg.Region,
		"endpoint":      dsc.cfg.Endpoint,
		"rootDirectory": dsc.cfg.RootDirectory,
	}
}

// Create is used to create our s3 datastore
func (dsc *DSConfig) Create(path string) (repo.Datastore, error) {
	d, err := s3.NewDatastore(dsc.cfg)
	if err != nil {
		return nil, err
	}
	if err := d.BucketExists(dsc.cfg.Bucket); err != nil {
		if err := d.CreateBucket(dsc.cfg.Bucket); err != nil {
			return nil, err
		}
	}
	return d, nil
}
