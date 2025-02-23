# api-gateway-config-supervisor
Syncs config files from Amazon S3, Google Drive, Dropbox, Amazon Cloud Drive reloading the gateway when there are any changes.

Table of Contents
=================

* [Status](#status)
* [Quick Start](#quick-start)
* [Dependencies](#dependencies)
* [Developer Guide](#developer-guide)

Status
======
This module is considered production ready. 

Quick Start
============

This module should be executed alongside the gateway in order to keep track of the configuration files and reload the gateway when there is a change.

```
api-gateway-config-supervisor \
        --reload-cmd="api-gateway -s reload" \
        --sync-folder=/etc/api-gateway \
        --sync-interval=10s \
        --sync-cmd="rclone sync s3-gw-config:api-gateway-config/ /etc/api-gateway -q" \
        --http-addr=127.0.0.1:8888
```

`sync-cmd` is executed each `sync-interval`. If there are changes to the files in `sync-folder` `reload-cmd` is executed.
`sync-folder` needs to exist before executing the command otherwise it exits with an error.
A web server is also started at `http-addr`; the gateway should check `http://<http-addr>/health-check` as part of its own regular `health-check`
so that in the unlikely event that this process dies the gateway appears unhealthy too.

In the initial design `rclone` was embedded into the program but b/c it wasn't straight forward to integrate it `sync-cmd` is used instead.
Using an external command for syncing is not that bad actually as it allows other cloud specific tools to come into play ( i.e `aws cli` )

### Using AWS-CLI
After installing `aws cli` the only change required to use this tool is to edit `sync-cmd` to something like:
```
--sync-cmd="aws s3 sync s3://api-gateway-config /etc/api-gateway"
```

```
api-gateway-config-supervisor \
        --reload-cmd="api-gateway -s reload" \
        --sync-folder=/etc/api-gateway \
        --sync-interval=10s \
        --sync-cmd="aws s3 sync s3://api-gateway-config /etc/api-gateway" \
        --http-addr=127.0.0.1:8888
```

Dependencies
============

* Golang
* https://gowalker.org/github.com/ncw/rclone

Developer guide
===============

Make sure you have go installed. On a Mac you can execute:
```
brew install go
```

To run the unit tests run:

```
make test
```

### Building a Docker image

Make sure you install docker first, opening a Docker terminal, then issue:

```
make docker
```

The `Dockerfile` is building a minimalistic Docker image installing go and its dependencies only to build the project, uninstalling them afterwards,
and only keeping statically built binaries. In addition it adds `rclone` ( +~14MB ) and `awscli` ( +~70MB ) for convenience but `awscli` is to be removed once `rclone` supports IAM Roles.

The container's entrypoint is `api-gateway-config-supervisor` so that its usage looks similar to the command without docker. For example:
```
docker run adobeapiplatform/api-gateway-config-supervisor:latest \
        --reload-cmd="api-gateway -s reload" \
        --sync-folder=/etc/api-gateway \
        --sync-interval=10s \
        --sync-cmd="aws s3 sync s3://api-gateway-config /etc/api-gateway" \
        --http-addr=127.0.0.1:8888
```