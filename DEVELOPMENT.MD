## notes

* [ipfs/plugins.md](https://github.com/ipfs/go-ipfs/blob/master/docs/plugins.md)
* [example plugin](https://github.com/ipfs/go-ipfs-example-plugin/)
* [ipfs/go-datastore](https://github.com/ipfs/go-datastore)
* [aws sdk for minio](https://docs.minio.io/docs/how-to-use-aws-sdk-for-go-with-minio-server.html)
* [s3 plugin pr](https://github.com/ipfs/go-ipfs/pull/5561)
* [s3 plugin code](https://github.com/ipfs/go-ipfs/blob/1526a4a7b2be3eb7c8dfba15dd64a8c8ebf021d6/plugin/plugins/s3ds/s3ds.go)

## development

* The easiest way to develop this is to run a local version of the storj v3 network see this [awesome tutorial](https://medium.com/@kleffew/getting-started-with-the-storj-v3-test-network-storj-sdk-c835d992cdd9) which covers everything you'll need to know.

## plugin

* `Create` is used to create the actual datastore
* `fsrepo.ConfigFromMap` is a function. not quite sure how it loads the configuration yet