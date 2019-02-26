package main

import (
	plugin "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/plugin"

	"github.com/RTradeLtd/storj-ipfs-ds-plugin/storj"
)

// Plugins is an exported list of plugins that will be loaded by go-ipfs.
var Plugins = []plugin.Plugin{
	&storj.SJPlugin{},
}

func main() {}
