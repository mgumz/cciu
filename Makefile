VERSION=$(shell cat VERSION)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH=$(shell git rev-parse HEAD)
TARGETS=linux.amd64 		\
		linux.arm64 		\
		linux.mips64 		\
		windows.amd64.exe 	\
		darwin.amd64 		\
		darwin.arm64
CONTAINER_IMAGE=quay.io/mgumz/cciu:$(VERSION)

BINARIES=$(foreach r,$(TARGETS),bin/cciu-$(VERSION).$(r))
RELEASES=$(subst windows.amd64.tar.gz,windows.amd64.zip,$(foreach r,$(subst .exe,,$(TARGETS)),releases/cciu-$(VERSION).$(r).tar.gz))

LDFLAGS=-trimpath 							\
		-ldflags "$(EXTRA_LDFLAGS) 			\
		-X main.versionString=$(VERSION)	\
		-X main.buildDate=$(BUILD_DATE)		\
		-X main.gitHash=$(GIT_HASH)"

.PHONY: toc
toc:
	@echo "list of targets:"
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | \
		awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | \
		sort | \
		egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | \
		awk '{ print " ", $$1 }'

.PHONY: binaries
binaries: $(BINARIES)
.PHONY: releases
releases: $(RELEASES)
	make $(RELEASES)
.PHONY: list-releases
list-releases:
	@echo $(RELEASES)|tr ' ' '\n'

.PHONY: clean
clean:
	rm -f $(BINARIES) $(RELEASES)

cciu: bin/cciu
.PHONY: cciu-small
cciu-small:
	make EXTRA_LDFLAGS="-s -w" bin/cciu
bin/cciu:
	go build $(LDFLAGS) -o $@ ./cmd/cciu

bin/cciu-$(VERSION).%:
	env GOARCH=$(subst .,,$(suffix $(subst .exe,,$@))) GOOS=$(subst .,,$(suffix $(basename $(subst .exe,,$@)))) CGO_ENABLED=0 \
	go build $(LDFLAGS) -o $@ ./cmd/cciu

releases/cciu-$(VERSION).%.zip: bin/cciu-$(VERSION).%.exe
	mkdir -p releases
	zip -9 -j -r $@ README.md LICENSE $<
releases/cciu-$(VERSION).%.tar.gz: bin/cciu-$(VERSION).%
	mkdir -p releases
	tar -cf $(basename $@) README.md LICENSE && \
		tar -rf $(basename $@) --strip-components 1 $< && \
		gzip -9 $(basename $@)


.PHONY: container-image
container-image: container/Containerfile
	docker build -f $< -t $(CONTAINER_IMAGE) .

.PHONY: deps-vendor
deps-vendor:
	go mod vendor
.PHONY: deps-cleanup
deps-cleanup:
	go mod tidy
.PHONY: deps-ls
deps-ls:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' all
.PHONY: deps-ls-updates
deps-ls-updates:
	go list -m -mod=readonly -f '{{if not .Indirect}}{{.}}{{end}}' -u all

.PHONY: reports


reports: report-golangci-lint
reports: report-vuln

report-golangci-lint:
	@echo '####################################################################'
	golangci-lint run ./cmd/... ./pkg/...

report-vuln:
	@echo '####################################################################'
	govulncheck ./cmd/... ./pkg/...

report-grype:
	@echo '####################################################################'
	grype .

fetch-report-tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest

fetch-report-tool-grype:
	go install github.com/anchore/grype@latest



.PHONY: cciu
.PHONY: bin/cciu
