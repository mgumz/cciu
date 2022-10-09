VERSION=0.8.1
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)
TARGETS=linux.amd64 linux.arm64 linux.mips64 windows.amd64.exe darwin.amd64 darwin.arm64
CONTAINER_IMAGE=quay.io/mgumz/cciu:$(VERSION)

LDFLAGS=$(EXTRA_LDFLAGS) -X main.versionString=$(VERSION) -X main.buildDate=$(BUILD_DATE) -X main.gitHash=$(GIT_HASH)
BINARIES=$(foreach r,$(TARGETS),bin/cciu-$(VERSION).$(r))
RELEASES=$(subst windows.amd64.tar.gz,windows.amd64.zip,$(foreach r,$(subst .exe,,$(TARGETS)),releases/cciu-$(VERSION).$(r).tar.gz))

toc:
	@echo "list of targets:"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | \
		awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | \
		sort | \
		egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | \
		awk '{ print " ", $$1 }'

binaries: $(BINARIES)
releases: $(RELEASES)
	make $(RELEASES)

clean:
	rm -f $(BINARIES) $(RELEASES)

cciu: bin/cciu
cciu-small:
	make EXTRA_LDFLAGS="-s -w" bin/cciu
bin/cciu:
	go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/cciu

bin/cciu-$(VERSION).%:
	env GOARCH=$(subst .,,$(suffix $(subst .exe,,$@))) GOOS=$(subst .,,$(suffix $(basename $(subst .exe,,$@)))) CGO_ENABLED=0 \
	go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/cciu

releases/cciu-$(VERSION).%.zip: bin/cciu-$(VERSION).%.exe
	mkdir -p releases
	zip -9 -j -r $@ README.md LICENSE $<
releases/cciu-$(VERSION).%.tar.gz: bin/cciu-$(VERSION).%
	mkdir -p releases
	tar -cf $(basename $@) README.md LICENSE && \
		tar -rf $(basename $@) --strip-components 1 $< && \
		gzip -9 $(basename $@)


container-image: docker/Dockerfile
	docker build -f docker/Dockerfile -t $(CONTAINER_IMAGE) .


deps-vendor:
	go mod vendor
deps-cleanup:
	go mod tidy
deps-ls:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' all
deps-ls-updates:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' -u all

report: report-cyclo report-lint report-staticcheck report-mispell report-ineffassign report-vet
report-cyclo:
	@echo '####################################################################'
	gocyclo ./cmd/cciu
report-mispell:
	@echo '####################################################################'
	misspell ./cmd/
report-lint:
	@echo '####################################################################'
	golint ./cmd/... ./internal/...
report-ineffassign:
	@echo '####################################################################'
	ineffassign ./cmd/... ./internal/...
report-vet:
	@echo '####################################################################'
	go vet ./cmd/... ./internal/...
report-staticcheck:
	@echo '####################################################################'
	staticcheck ./cmd/... ./internal/...

fetch-report-tools:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
 

.PHONY: cciu bin/cciu
