# storj-ipfs-ds-plugin


This plugin is used to enable S3 datastore usage through STORJ. The immediate benefit to this is an immense amount of data durability, and redundancy not present with current IPFS node solutions. Typically to achieve this, you'll need to run many hosts, copies of the data, and hardware or software level data protection methods like RAID or ZFS. This is currently an experimental production, and it not advised to use this for production data. There are a few known bugs, and probably some hidden ones.

Due to using the plugin system which is heavily dependent on gx, until further notice it's best to use the binaries in the `bin` folder. If you want to rebuild the `ipfs` binary you'll wnat to go into the `vendor` folder and build directory instead of using a prebuilt binary from a separate code-base.

## warnings

* When uploading data, sometimes it may appear as if the upload has stalled. This is currently due to the Segments being uploaded sequentially, while the pieces (fragments of the segments) are uploaded in parallel.

* sometimes after restarting your IPFS node, it may appear as if data has been lost by running `ipfs pin ls` and nothing showing up. This shouldn't always happen, but if it does you can verify data is stored by shutting down your node, and running `ipfs pin ls` a second time which should show the data. Obviously this is not viable in production, so this will be ironed out before recommending usage of this in production

## demo 

short video of daemon operation:
https://gateway.temporal.cloud/ipfs/QmeFisZdZuHmnwaXEUBCaMJmoHQLLPn3DJfNiYwdCug5iG

## configuration

The following is an example IPFS configuration to use this plugin.

```json
~/.ipfs/config example
{
    // ...
    "Datastore": {
    "StorageMax": "10GB",
    "StorageGCWatermark": 90,
    "GCPeriod": "1h",
    "Spec": {
      "mounts": [
        {
          "child": {
            "accessKey": "...",
            "secretKey": "...",
            "bucket": "go-ipfs-storj",
            "region": "us-east-1",
            "endpoint": "127.0.0.1:9000",
            "rootDirectory": "",
            "workers" "100",
            "type": "storj"
          },
          "mountpoint": "/blocks",
          "name": "storj",
          "type": "log"
        },
        {
          "child": {
            "compression": "none",
            "path": "datastore",
            "type": "levelds"
          },
          "mountpoint": "/",
          "prefix": "leveldb.datastore",
          "type": "measure"
        }
      ],
      "type": "mount"
    },
    "HashOnRead": false,
    "BloomFilterSize": 10000000
  }, // ...
}
```
