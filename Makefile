GOPATH ?= `pwd`
GOBIN ?= `pwd`/bin
GOOS ?= $(`uname -a | awk '{print tolower($1)}'`)
.PHONY: setup install test format static docker docker-ssh

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
