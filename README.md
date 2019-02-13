# storj-ipfs-ds-plugin

go-ipfs plugin to use STORJ as the storage backend

## what works and doesn't work

* currently all that works is basic daemon running
* any datastore operations (put for example) do not work yet as there is no Batch operation system yet

short video of daemon operation:
https://gateway.temporal.cloud/ipfs/QmeFisZdZuHmnwaXEUBCaMJmoHQLLPn3DJfNiYwdCug5iG

## notes

* [ipfs/plugins.md](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md)
* [example plugin](https://github.com/ipfs/go-ipfs-example-plugin/)
* [ipfs/go-datastore](https://github.com/ipfs/go-datastore)
* [aws sdk for minio](https://docs.minio.io/docs/how-to-use-aws-sdk-for-go-with-minio-server.html)
* [s3 plugin pr](https://github.com/ipfs/go-ipfs/pull/5561)
* [s3 plugin code](https://github.com/ipfs/go-ipfs/blob/1526a4a7b2be3eb7c8dfba15dd64a8c8ebf021d6/plugin/plugins/s3ds/s3ds.go)

## plugin

* `Create` is used to create the actual datastore
* `fsrepo.ConfigFromMap` is a function. not quite sure how it loads the configuration yet

## configuration

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
