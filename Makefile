SOURCEDIR := .
OUTDIR := bin
DATE := $(shell date +%Y-%m-%d-%H-%M-%S)
GITHASH := $(shell git show -s --format=%h)

.PHONY: build
build:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -X \"main.Version=$(PACKAGE_VERSION)-$(GITHASH)\" -X \"main.BuildDate=$(DATE)\" -X \"main.GoVersion=`go version|cut -d ' ' -f 3`\"" -o $(OUTDIR)/shync_darwin_amd64 cmd/shync/main.go
	# GOOS=linux GOARCH=amd64 go build -ldflags="-s -X \"main.Version=$(PACKAGE_VERSION)-$(GITHASH)\" -X \"main.BuildDate=$(DATE)\" -X \"main.GoVersion=`go version|cut -d ' ' -f 3`\"" -o $(OUTDIR)/shync_linux_amd64 cmd/shync/main.go

run:
	SHYNC_LOGLEVEL=debug go run cmd/shync/main.go download .

.PHONY: test
test:
	SHYNC_LOGLEVEL=debug go test -v -race ./... -count=1

clean:
	rm -f bin/*
