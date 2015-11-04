GOPATH ?= `pwd`
GOBIN ?= `pwd`/bin

install:
	@go version
	export GO15VENDOREXPERIMENT=1
#	GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install -v ./...
	GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(GOPATH)/bin/godep go install -v ./...

format:
	gofmt -e -w ./

docker:
	docker build -t adobeapiplatform/api-gateway-config-supervisor .