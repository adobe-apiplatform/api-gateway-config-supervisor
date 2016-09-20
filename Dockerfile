# api-gateway-config-supervisor
#
# VERSION               1.9.3.1
#
# From https://hub.docker.com/_/alpine/
#
FROM alpine:latest

ENV GOPATH /tmp/go:/usr/lib/go
ENV GOBIN  /usr/lib/go/bin
ENV PATH   $PATH:/usr/lib/go/bin


RUN mkdir -p /tmp/go/src/github.com/adobe-apiplatform/api-gateway-config-supervisor
ADD . /tmp/go/src/github.com/adobe-apiplatform/api-gateway-config-supervisor

RUN echo " installing aws-cli ..." \
    && echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && apk update \
    && apk add python \
    && apk add py-pip \
    && pip install --upgrade pip \
    && pip install awscli \
    && apk del py-pip \

    && echo " building local project ... " \
    && apk add make git gcc libc-dev go \
    && cd /tmp/go/src/github.com/adobe-apiplatform/api-gateway-config-supervisor \
    && make setup \
    && godep  go test \
    && godep  go build -ldflags "-s" -a -installsuffix cgo -o api-gateway-config-supervisor ./ \
    && cp /tmp/go/src/github.com/adobe-apiplatform/api-gateway-config-supervisor/api-gateway-config-supervisor /usr/lib/go/bin \

    && echo "installing rclone sync ... " \
    && go get github.com/ncw/rclone \

    && echo " cleaning up ... " \
    && rm -rf /usr/lib/go/bin/src \
       	      /tmp/go \
    	      /tmp/api-gateway-config-supervisor* \
    	      /usr/lib/go/bin/pkg/ \
    	      /usr/lib/go/bin/godep \
    && apk del make git gcc libc-dev go \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["api-gateway-config-supervisor"]
