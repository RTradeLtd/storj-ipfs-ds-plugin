package node

import (
	"fmt"
	"os"
	"testing"
	"time"
)

/*
   MINIO_ACCESS_KEY: "C03T49S17RP0APEZDK6M"
   MINIO_SECRET_KEY: "q4I9t2MN/6bAgLkbF6uyS7jtQrXuNARcyrm2vvNA"
*/

const (
	accessKey = "C03T49S17RP0APEZDK6M"
	secretKey = "q4I9t2MN/6bAgLkbF6uyS7jtQrXuNARcyrm2vvNA"
)

var (
	homeDir        = os.Getenv("HOME")
	ipfsDir        = homeDir + "/.ipfs"
	ipfsConfigFile = ipfsDir + "/config"
	repoDir        = ipfsDir + "/blocks"
)

func TestNode(t *testing.T) {
	node, err := NewNode(accessKey, secretKey, ipfsConfigFile, repoDir)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("sleeping for 10 seconds")
	time.Sleep(time.Second * 10)
	if err := node.Close(); err != nil {
		t.Fatal(err)
	}
}
