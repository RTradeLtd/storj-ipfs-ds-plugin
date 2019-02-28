
# Storj-IPFS-DS-Plugin (Alpha)

***DO NOT USE FOR IMPORTANT DATA***

This repository contains code to facilitate running an IPFS node using the Storj network as the datastorage backend. It does this leveraging a good chunk of code that was borrowed from [go-ds-s3](https://github.com/ipfs/go-ds-s3).

The immediate benefit to a system like this is a mass amount of data protection (redundancy, decentralization, and durability) not present with traditional IPFS nodes. Traditional IPFS nodes require many copies of the data to spread it amongst multiple hosts and to achieve redundancy measures. However, when combined with Storj a single IPFS node is atuomatically distributing data across a vast number of machines accomplishing:

1) Data distribution
2) Data decentralization
3) Data redundancy
4) *Native data encryption at rest* <sup>1</sup>

1: Data is stored within the Storj network encrypted, however this does not imply that the data is encrypted to the IPFS network. Simply sharing the content hash to someone is enough for them to be able to pull the data and decrypt it.

## Contents

* `node` folder is a work-in-progress purpose-built IPFS node designed to use the Storj network as the data storage backend.
* `s3` folder is a modified version of [go-ds-s3](https://github.com/ipfs/go-ds-s3)
* `storj` folder is where the actual plugin lives
* `log` folder is a wrapper around the `zap` package by uber
* `fsrepo` is a hack of `go-ipfs/repo/fsrepo` to allow for a storj configuration profile during `ipfs init`
* `go-ipfs-repo` is a hack of `go-ipfs-config` to allow for a storj configuration profile during `ipfs init`

## Install

***Note 1: the build process assumes that the IPFS repository will be placed in `$HOME/.ipfs`, if you want to point to another path you'll need to edit the make file.***

***Note 2: This must be worked on in an appropriate Golang environment designed to work with `dep` (ie, package must be downloaded to `GOPATH/src/github.com/RTradeLtd/storj-ipfs-ds-plugin)***

To build from, and install from source on a machine which already has an initialized IPFS instance, you can run the following command.

```shell
git clone https://github.com/RTradeLtd/storj-ipfs-ds-plugin.git
cd storj-ipfs-ds-plugin
make install
```

If you aren't running this on a machine that already has an initialized IPFS node, running the following commands will handle the initialization process

```shell
git clone https://github.com/RTradeLtd/storj-ipfs-ds-plugin.git
cd storj-ipfs-ds-plugin
make first-install
```


Running either of the above `make` commands will create a folder `build` and place both the IPFS binary, and Storj plugin inside.

When running `ipfs init --profile=storj` (which is what `make first-install`) does, you can automatically configure your s3 access key, and s3 secret key with the following environment variables, otherwise you'll have to manually edit `$IPFS_PATH/config`

```shell
STORJ_ACCESS_KEY=...
STORJ_SECRET_KEY=...
```

## Developing

This uses the `0.4.18` release of the IPFS code base, with a few hacks to `fsrepo` and `go-ipfs-config` get the ipfs daemon to recognize `storj` as a valid configuration profile. It is recommended that you develop against the `storj-sim` network, for which an excellent tutorial can be seen on [medium](https://medium.com/@kleffew/getting-started-with-the-storj-v3-test-network-storj-sdk-c835d992cdd9)

After making any changes to `fsrepo` run `make fsrepo` to update the version used by the bundled go-ipfs dependency

After making any changes to `go-ipfs-config` run `make go-ipfs-config` to update the version used by the bundled go-ipfs dependency

### Dependencies

All dependencies needed are shipped with this repository. It uses go-ipfs 0.4.18, with all IPFS dependencies managed by `gx`. Please see `package.json` for a detailed list of IPFS dependencies.

None IPFS dependencies are managed by `dep`

This has only been tested with go1.11+

## Warnings

Being alpha level software there are a ton of gotchas please read the following in detail

### Superficial Data Loss

This is still very much experimental software and lots of issues are present. The most notable issue is that when storing data into S3 via Storj, anytime you restart a node it appears as if "data is lost".  Data is not actually lost, and it is still stored in the S3 interface, but running a command such as `ipfs pin ls` shows that nothing is stored. Additionally running `ipfs repo stat` shows that barely any space is being taken up, and that there are no objects

There are a few potential causes for this:

1) the `$IPFS_PATH/blocks` and `$IPFS_PATH/datastore` not containing anything
2) some of the encryption/path encryption that storj performs
3) Data distribution due to reed-solomon
4) ????

### Inconsistencies

There are numerous data inconsistencies (albeit superficial). For example note the following output from `ipfs repo stat`:

```shell
solidity@dark:~/go/src/github.com/RTradeLtd/storj-ipfs-ds-plugin/build$ ./ipfs repo stat
NumObjects: 0
RepoSize:   28881
StorageMax: 10000000000
RepoPath:   /home/solidity/.ipfs
Version:    fs-repo@7
```

But now note the output from `ipfs pin ls | wc -l`:

```shell
solidity@dark:~/go/src/github.com/RTradeLtd/storj-ipfs-ds-plugin/build$ ./ipfs pin ls | wc -l
162
```

### Performance

Performance when running in a local development environment (ie, `storj-sim`) is noticeably slower than when using other IPFS datastores, and currently when using a remote storj network (ie, not `storj-sim`) the performance tanks even more.

## Configuration

### IPFS

If you need to manually configure `$IPFS_PATH/plugins` you can do so with the following snippet (note, you'll need to change settings as desired):

```json
  "Datastore": {
    "StorageMax": "10GB",
    "StorageGCWatermark": 90,
    "GCPeriod": "1h",
    "Spec": {
      "mounts": [
        {
          "child": {
            "accessKey": "45qhD6QdHjrNJ65vua16ZoRwEtTV",
            "secretKey": "4MFwsyGXCuuND15Wt37P7MH8HvVV",
            "bucket": "go-ipfs-storj-5",
            "region": "us-east-1",
            "endpoint": "http://127.0.0.1:9000",
            "rootDirectory": "",
            "type": "storj",
            "logPath": ""
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
    "BloomFilterSize": 0
  },
```

You will then need to update `$IPFS_PATH/datastore_spec` to reflect the above:

```json
{"mounts":[{"bucket":"go-ipfs-storj-5","endpoint":"http://127.0.0.1:9000","mountpoint":"/blocks","region":"us-east-1","rootDirectory":""},{"mountpoint":"/","path":"datastore","type":"levelds"}],"type":"mount"}
```

### Storj

When using a remote satellite,, you'll probably want to update your rs (reed-solomon) settings of the Storj minio gateway to something like

```yaml
# the largest amount of pieces to encode to. n.
rs.max-threshold: 50
# the minimum pieces required to recover a segment. k.
rs.min-threshold: 20
# the minimum safe pieces before a repair is triggered. m.
rs.repair-threshold: 25
# the desired total pieces for a segment. o.
rs.success-threshold: 40
``