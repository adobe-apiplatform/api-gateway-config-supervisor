GOPATH ?= `pwd`
GOBIN ?= `pwd`/bin
GOOS ?= $(`uname -a | awk '{print tolower($1)}'`)

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

DOCKER_TEMP := $(shell mktemp -u ./.Dockerfile.XXXXXXXXXX)
.DELETE_ON_ERROR: ${DOCKER_TEMP}
${DOCKER_TEMP}:
ifdef DOCKER_BUILD_DEBUG
	cp Dockerfile ${DOCKER_TEMP}
else
	perl -0777 -pe 's/#.*?\n+//g; s/\n+/\n/g; s{(\nRUN )}{++$$count > 1 ? "\\\n&& " : $$1}ge;' <Dockerfile >${DOCKER_TEMP}
endif
	docker build -t adobeapiplatform/api-gateway-config-supervisor -f ${DOCKER_TEMP} .
	rm -f ${DOCKER_TEMP}

.PHONY: docker
docker: ${DOCKER_TEMP}

.PHONY: docker-ssh
docker-ssh:
	docker run -ti --entrypoint='/bin/sh' adobeapiplatform/api-gateway-config-supervisor:latest
