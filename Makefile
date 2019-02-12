IPFSVERSION=v0.4.18


.PHONY: testenv
testenv:
	(cd testenv ; make testenv)

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
	dep ensure -v -update

	cp -r /tmp/gx vendor/gx
	
	# Remove problematic dependencies
	find . -name test-vectors -type d -exec rm -r {} +
	@echo "===================          done           ==================="


# install gx related management dependencies
.PHONY: gx
gx:
	go get -u -v github.com/whyrusleeping/gx