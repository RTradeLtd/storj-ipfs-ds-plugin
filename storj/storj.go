package storj

import (
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/plugin"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo"
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
