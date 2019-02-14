# notes

So to be a valid repo, it looks like we need to satisfy the following

```Golang
// FSRepo represents an IPFS FileSystem Repo. It is safe for use by multiple
// callers.
type FSRepo struct {
	// has Close been called already
	closed bool
	// path is the file-system path
	path string
	// lockfile is the file system lock to prevent others from opening
	// the same fsrepo path concurrently
	lockfile io.Closer
	config   *config.Config
	ds       repo.Datastore
	keystore keystore.Keystore
	filemgr  *filestore.FileManager
}
```

The s3 datastore will satisfy the `repo.Datastore` of the above struct, `config.Config` is handled as well. Not sure about the other parts


## creating a custom repo

```
* use fsrepo
* this link is how ipfs does initializes itthttps://github.com/ipfs/go-ipfs/blob/402af03196c4fa9fff62e9942992fe76963c082b/cmd/ipfs/init.go#L134-L178
* need to create a go-ipfs config, call fsrepo.Init
* you can then open it with fsrepo.Open
```