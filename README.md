
# Storj-IPFS-DS-Plugin (Alpha)

## DO NOT USE FOR IMPORTANT DATA

This repository contains code to facilitate running an IPFS node using the Storj network as the datastorage backend. It does this leveraging a good chunk of code that was borrowed from [go-ds-s3](https://github.com/ipfs/go-ds-s3).

The immediate benefit to a system like this is a mass amount of data protection (redundancy, decentralization, and durability) not present with traditional IPFS nodes. Traditional IPFS nodes require many copies of the data to spread it amongst multiple hosts and to achieve redundancy measures. However, when combined with STORJ a single IPFS node is atuomatically distributing data across a vast number of machines accomplishing:

1) Data distribution
2) Data decentralization
3) Data redundancy
4) *Native data encryption* (data is store within the Storj network encrypted)

To see a short,but old video of daemon operation see [here](https://gateway.temporal.cloud/ipfs/QmeFisZdZuHmnwaXEUBCaMJmoHQLLPn3DJfNiYwdCug5iG)

## Install

***Note: the build process assumes that the IPFS repository will be placed in `$HOME/.ipfs`, if you want to point to another path you'll need to edit the make file.***

To build from, and install from source on a machine which already has an initialized IPFS instance, you can run the following command

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

When running `ipfs init --profile=storj` (which is what `make first-install`) does, you can automatically configure your s3 access key, and s3 secret key with the following environment variables, otherwise you'll have to manually edit `$IPFS_PATH/config`

```shell
STORJ_ACCESS_KEY=...
STORJ_SECRET_KEY=...
```


## Developing

This uses the `0.4.18` release of the IPFS code base, with a few hacks to get the ipfs daemon to recognize `storj` as a valid configuration profile:

After making any changes to `fsrepo` run the following commands:

```shell
rm -rf ./vendor/gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo
cp -r fsrepo ./vendor/gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo
```

After making any changes to `go-ipfs-config` run the following commands:

```shell
rm -rf ./vendor/gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config
cp -r go-ipfs-config ./vendor/gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config
```

## Contents

* `node` folder is a work-in-progress purpose-built IPFS node designed to use the Storj network as the data storage backend.
* `s3` folder is a modified version of [go-ds-s3](https://github.com/ipfs/go-ds-s3) 
* `storj` folder is where the actual plugin lives
* `log` folder is a wrapper around the `zap` package by uber

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

## Dependencies

All dependencies needed are shipped with this repository. It uses go-ipfs 0.4.18, with all IPFS dependencies managed by `gx`. Please see `package.json` for a detailed list of IPFS deependencies.

None IPFS dependencies are managed by `dep`

This has only been tested with go1.11+

## Installation

The install process is simple, running `make install` will do the following (in order):

1) build the ipfs daemon from bundled dependencies
2) build the plugin
3) install plugin into `~/.ipfs/plugins`

You can then use the `ipfs` binary included in the `build` folder

## Configuration

You will need to update `$IPFS_PATH/config` as there is not yet functionality to do this for you. The following is an example configuration, which you will want to change depending on your environment:

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

You will then need to update `$IPFS_PATH/datastore_spec` to reflect the boave:

```json
{"mounts":[{"bucket":"go-ipfs-storj-5","endpoint":"http://127.0.0.1:9000","mountpoint":"/blocks","region":"us-east-1","rootDirectory":""},{"mountpoint":"/","path":"datastore","type":"levelds"}],"type":"mount"}
```

When using a remote satellite, you'll want to update the rs (reed-solomon) settings of your Storj IPFS Gateway to something like:

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