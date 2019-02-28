GXIPFSVERSION=QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r
IPFSCMDBUILDPATH=vendor/gx/ipfs/$(GXIPFSVERSION)/go-ipfs/cmd/ipfs
REPOROOT=$(shell pwd)
IPFSPATH=$(HOME)/.ipfs

# starts minio docker container for testing
.PHONY: testenv
testenv:
	(cd testenv ; make minio)

# cleans up test environment
.PHONY: stop-testenv
stop-testenv:
	(cd testenv ; make clean)
	
# Rebuild vendored dependencies
.PHONY: vendor
vendor:
	@echo "=================== generating dependencies ==================="
	cp -r vendor/gx /tmp/gx
	# Nuke vendor directory
	rm -rf vendor

	# Update standard dependencies
	dep ensure -v $(DEPFLAGS)

	cp -r /tmp/gx vendor/gx
	rm -rf vendor/gx/gx
	# Remove problematic dependencies
	find . -name test-vectors -type d -exec rm -r {} +
	@echo "===================          done           ==================="

# used to update the fsrepo package used by ipfs
.PHONY: fsrepo
fsrepo:
	rm -rf ./vendor/gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo
	cp -r fsrepo ./vendor/gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/repo/fsrepo

# used to update the go-ipfs-config package used by ipfs
.PHONY: go-ipfs-config
go-ipfs-config:
	rm -rf ./vendor/gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config
	cp -r go-ipfs-config ./vendor/gx/ipfs/QmPEpj17FDRpc7K1aArKZp3RsHtzRMKykeK9GVgn4WQGPR/go-ipfs-config

# install is used if you have already installed IPFS on your workstation
.PHONY: install
install: build-plugin install-plugin

# first-install is used if you have never installed IPFS on your workstation
.PHONY: first-install
first-install: build-plugin init install-plugin

# installs the plugin
.PHONY: install-plugin
install-plugin:
	rm -rf $(IPFSPATH)/plugins
	mkdir $(IPFSPATH)/plugins
	install -Dm700 build/storj-ipfs-ds-plugin.go $(IPFSPATH)/plugins

# builds the actual plugin and ipfs node
.PHONY: build-plugin
build-plugin:
	mkdir $(REPOROOT)/build
	(cd $(IPFSCMDBUILDPATH) ; go build ; cp ipfs $(REPOROOT)/build)
	(go build -o build/storj-ipfs-ds-plugin.go --buildmode=plugin ; chmod a+x build/storj-ipfs-ds-plugin.go)

# initializes an ipfs node, after first having built the plugin
.PHONY: init
init:
	(cd $(REPOROOT)/build ; ./ipfs init --profile=storj)

# cleans up build files
.PHONY: clean
clean:
	rm -rf build

# install gx related management dependencies
.PHONY: gx
gx:
	go get -u -v github.com/whyrusleeping/gx


# list make targets
.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs