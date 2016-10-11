GOPATH ?= `pwd`
GOBIN ?= `pwd`/bin
GOOS ?= $(`uname -a | awk '{print tolower($1)}'`)
.PHONY: setup install test format static docker docker-ssh

setup:
	go get github.com/tools/godep

install:
	@go version
	export GO15VENDOREXPERIMENT=1
#	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install -v ./...
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOPATH)/bin/godep go install -v ./...

test:
	go test

format:
	gofmt -e -w ./

static:
#	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o $(GOBIN)/api-gateway-config-supervisor-static ./
	CGO_ENABLED=0 go build -ldflags "-s" -a -installsuffix cgo -o $(GOBIN)/api-gateway-config-supervisor-static ./

docker:
	docker build -t adobeapiplatform/api-gateway-config-supervisor .

docker-ssh:
	docker run -ti --entrypoint='/bin/sh' adobeapiplatform/api-gateway-config-supervisor:latest
