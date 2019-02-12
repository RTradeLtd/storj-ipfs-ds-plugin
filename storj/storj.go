package storj

import (
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
	return nil
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
