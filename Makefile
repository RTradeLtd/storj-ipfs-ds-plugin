GXIPFSVERSION=QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r
IPFSCMDBUILDPATH=vendor/gx/ipfs/$(GXIPFSVERSION)/go-ipfs/cmd/ipfs
REPOROOT=$(GOPATH)/src/github.com/RTradeLtd/storj-ipfs-ds-plugin
IPFSPATH=$(HOME)/.ipfs

.PHONY: testenv
testenv:
	(cd testenv ; make minio)

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

# install is used if you have already installed IPFS on your workstation
.PHONY: install
install: build-plugin install-plugin

# first-install is used if you have never installed IPFS on your workstation
.PHONY: first-install
first-install: build-plugin init install-plugin

.PHONY: install-plugin
install-plugin:
	rm -rf $(IPFSPATH)/plugins
	mkdir $(IPFSPATH)/plugins
	install -Dm700 build/storj-ipfs-ds-plugin.go $(IPFSPATH)/plugins

.PHONY: build-plugin
build-plugin:
	mkdir $(REPOROOT)/build
	(cd $(IPFSCMDBUILDPATH) ; go build ; cp ipfs $(REPOROOT)/build)
	(go build -o build/storj-ipfs-ds-plugin.go --buildmode=plugin ; chmod a+x build/storj-ipfs-ds-plugin.go)

.PHONY: init
init:
	(cd $(REPOROOT)/build ; ./ipfs init)

.PHONY: clean
clean:
	rm -rf build

# install gx related management dependencies
.PHONY: gx
gx:
	go get -u -v github.com/whyrusleeping/gx