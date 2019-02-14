package node

import (
	"context"
	"encoding/base64"
	config "gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config"
	ci "gx/ipfs/QmPvyPwuCgJ7pDmrKDxRtsScJgBaM5h4EpRL2qQJsmXf4n/go-libp2p-crypto"
	peer "gx/ipfs/QmTRhk7cgjUf2gfQ3p2M9KPECNZEW9XUrmHcFCgog4cPgB/go-libp2p-peer"
	core "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core"
	repo "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo"
	ds "gx/ipfs/QmaRb5yNXKonhbkpNxNawoydk4N6es6b4fPj19sjEKsh5D/go-datastore"

	"github.com/RTradeLtd/storj-ipfs-ds-plugin/s3"
)

// SNode is our custom built storj ipfs node
type SNode struct {
	h *core.IpfsNode
	d ds.Datastore
}

// NewNode generates our storj backed ipfs node
func NewNode(accessKey, secretKey string) (*SNode, error) {
	datastoreConfig := s3.NewConfig(accessKey, secretKey)
	datastore, err := s3.NewDatastore(datastoreConfig)
	if err != nil {
		return nil, err
	}
	// setup datastore if needed
	if err := datastore.BucketExists(datastore.Config.Bucket); err != nil {
		if err := datastore.CreateBucket(datastore.Config.Bucket); err != nil {
			return nil, err
		}
	}
	pk, _, err := ci.GenerateKeyPair(ci.Ed25519, 258)
	if err != nil {
		return nil, err
	}
	pid, err := peer.IDFromPrivateKey(pk)
	if err != nil {
		return nil, err
	}
	pkBytes, err := pk.Bytes()
	if err != nil {
		return nil, err
	}
	nodeConfig := config.Config{}
	nodeConfig.Bootstrap = config.DefaultBootstrapAddresses
	nodeConfig.Identity.PeerID = pid.Pretty()
	nodeConfig.Identity.PrivKey = base64.StdEncoding.EncodeToString(pkBytes)
	// WARNING: repo.Mock is not threadsafe we will need to move from this
	mockRepo := repo.Mock{
		C: nodeConfig,
		D: datastore,
	}
	// create the node
	host, err := core.NewNode(context.Background(), &core.BuildCfg{
		Online:    true,
		Permanent: true,
		Repo:      &mockRepo,
		ExtraOpts: map[string]bool{
			"ipnsps": true,
		},
	})
	if err != nil {
		return nil, err
	}
	return &SNode{
		h: host,
		d: datastore,
	}, nil
}

// Close is used to shutdown our snode
func (sn *SNode) Close() error {
	return sn.h.Close()
}
