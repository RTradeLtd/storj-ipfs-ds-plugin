# storj-ipfs-ds-plugin

go-ipfs plugin to use STORJ as the storage backend

## notes

[ipfs/plugins.md](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md)
[example plugin](https://github.com/ipfs/go-ipfs-example-plugin/)
[ipfs/go-datastore](https://github.com/ipfs/go-datastore)
[aws sdk for minio](https://docs.minio.io/docs/how-to-use-aws-sdk-for-go-with-minio-server.html)
[s3 plugin wip](https://github.com/ipfs/go-ipfs/pull/5561)

A datastore must also have two functions `DiskSpec` and `Create`.

* `Create` is used to create the actual datastore