# api-gateway-config-supervisor
#
# VERSION               1.9.3.1
#
# From https://hub.docker.com/_/alpine/
#
FROM alpine:latest

ENV GOPATH=/tmp/go \
    GOBIN=/usr/lib/go/bin \
    PATH=$PATH:/usr/lib/go/bin

COPY . ${GOPATH}/src/github.com/adobe-apiplatform/api-gateway-config-supervisor

RUN echo " installing aws-cli ..." \
    && apk update \
    && apk add python \
    && apk add py-pip \
    && pip install --upgrade pip \
    && pip install awscli \
    && apk del py-pip \

    && echo " building local project ... " \
    && apk add --virtual .build-deps make git gcc libc-dev go \
    && cd ${GOPATH}/src/github.com/adobe-apiplatform/api-gateway-config-supervisor \
    && make setup \
    && godep  go test \
    && godep  go build -ldflags "-s" -a -installsuffix cgo -o api-gateway-config-supervisor ./ \
    && cp ${GOPATH}/src/github.com/adobe-apiplatform/api-gateway-config-supervisor/api-gateway-config-supervisor ${GOBIN} \

    && echo "installing rclone sync ... " \
    && go get github.com/ncw/rclone \

    && echo " cleaning up ... " \
    && rm -rf ${GOBIN}/src \
       	      ${GOPATH} \
    	      /tmp/api-gateway-config-supervisor* \
    	      ${GOBIN}/pkg/ \
    	      ${GOBIN}/godep \
    && apk del .build-deps \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["api-gateway-config-supervisor"]
