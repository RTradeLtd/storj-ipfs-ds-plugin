package node

import (
	"context"
	"encoding/json"
	config "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	core "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"
	"gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo"
	"os"

	sPlugin "github.com/RTradeLtd/storj-ipfs-ds-plugin/s3ds/plugin"
)

// SNode is our custom built storj ipfs node
type SNode struct {
	h *core.IpfsNode
}

// NewNode generates our storj backed ipfs node
func NewNode(ipfsConfigPath, repoPath string) (*SNode, error) {
	// genarate an empty node config
	nodeConfig := config.Config{}
	// open node config at provided path
	f, err := os.Open(ipfsConfigPath)
	if err != nil {
		return nil, err
	}
	// decode the opened config file into our config struct
	if err := json.NewDecoder(f).Decode(&nodeConfig); err != nil {
		return nil, err
	}
	// add our custom config handler
	if err := fsrepo.AddDatastoreConfigHandler("s3ds", sPlugin.DatastoreConfig); err != nil {
		return nil, err
	}
	// init our repo configuration
	if err := fsrepo.Init(repoPath, &nodeConfig); err != nil {
		return nil, err
	}
	// open the repo configuration
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}
	// create the node
	host, err := core.NewNode(context.Background(), &core.BuildCfg{
		Online:    true,
		Permanent: true,
		Repo:      repo,
		ExtraOpts: map[string]bool{
			"ipnsps": true,
		},
	})
	if err != nil {
		return nil, err
	}
	return &SNode{
		h: host,
	}, nil
}

// Close is used to shutdown our snode
func (sn *SNode) Close() error {
	return sn.h.Close()
}
