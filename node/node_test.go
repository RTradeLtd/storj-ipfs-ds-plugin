package node

import (
	"fmt"
	gocid "gx/ipfs/QmPSQnBKM9g7BaUcZCvswUJVscQ1ipjmwxN5PXCjkp9EQ7/go-cid"
	"os"
	"testing"
	"time"
)

/*
   MINIO_ACCESS_KEY: "C03T49S17RP0APEZDK6M"
   MINIO_SECRET_KEY: "q4I9t2MN/6bAgLkbF6uyS7jtQrXuNARcyrm2vvNA"
   	pin "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/pin"

*/

const (
	accessKey = "C03T49S17RP0APEZDK6M"
	secretKey = "q4I9t2MN/6bAgLkbF6uyS7jtQrXuNARcyrm2vvNA"
	testCID   = "QmS4ustL54uo8FzR9455qaxZwuMiUhyvMcX9Ba8nUH4uVv"
)

var (
	homeDir        = os.Getenv("HOME")
	ipfsDir        = homeDir + "/.ipfs"
	ipfsConfigFile = ipfsDir + "/config"
)

func TestNode(t *testing.T) {
	node, err := NewNode(accessKey, secretKey, ipfsConfigFile, ipfsDir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		fmt.Println("shutting down in 10 seconds")
		time.Sleep(time.Second * 10)
		fmt.Println(node.Close())
	}()
	cid, err := gocid.Decode(testCID)
	if err != nil {
		t.Fatal(err)
	}
	//node.h.Pinning.PinWithMode(cid, pin.Recursive)
	reason, pinned, err := node.h.Pinning.IsPinned(cid)
	if err != nil {
		t.Fatal(err)
	}
	if !pinned {
		fmt.Println("not pinned because of ", reason)
	} else {
		fmt.Println("data is pinned")
	}
}
