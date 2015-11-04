# api-gateway-config-supervisor
#
# VERSION               1.9.3.1
#
# From https://hub.docker.com/_/alpine/
#
FROM alpine:latest

RUN apk update \
    && apk add curl python

ENV GOLANG_VERSION 1.5.1
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz
ENV GOLANG_DOWNLOAD_SHA1 0df564746d105f4180c2b576a1553ebca9d9a124

RUN mkdir -p /tmp/api-gateway/

RUN echo "Installing Go ..." \
    && mkdir -p /tmp/api-gateway/ \
    && curl -L "$GOLANG_DOWNLOAD_URL" -o /tmp/api-gateway/golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA1  /tmp/api-gateway/golang.tar.gz" | sha1sum -c - \
    && mkdir -p /usr/local/src \
    && tar -C /usr/local -xzf /tmp/api-gateway/golang.tar.gz \
    && rm /tmp/api-gateway/golang.tar.gz