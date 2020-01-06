JOBDATE		?= $(shell date -u +%Y-%m-%dT%H%M%SZ)
GIT_REVISION	= $(shell git rev-parse --short HEAD)
VERSION		?= $(shell git describe --tags --abbrev=0)
GOOS ?= linux

test:
	go get github.com/mfridman/tparse	
	go test -json  -v `go list ./... | egrep -v /tests/` -cover | tparse -all -smallscreen
