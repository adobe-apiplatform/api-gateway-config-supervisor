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
This module is experimental at the moment.

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
A web server is also started at `http-addr`; the gateway should check `http://<http-addr>/health-check` as part of its own regular `health-check`
so that in the unlikely event that this process dies the gateway appears unhealthy too.

In the initial design `rclone` was embedded into the program but b/c it wasn't straight forward to integrate it `sync-cmd` is used instead.
Using an external command for syncing is not that bad actually as it allows other cloud specific tools to come into play ( i.e `aws cli` )

Dependencies
============

* https://github.com/tools/godep
* https://gowalker.org/github.com/ncw/rclone
* gopkg.in/fsnotify.v1
* TBD

Developer guide
===============

Make sure you have go 1.5.1+ installed. On a Mac you can execute
```
brew install go
```

The go dependencies have been already added using `godeps` in the existing GitRepo in order to achieve repeatable builds and isolated envs.
To build the project you can run:

```
make install
```

#### Building a Docker image

Make sure you install docker first, opening a Docker terminal, then issue:

```
make docker
```

The `Dockerfile` is building a minimalistic Docker image installing go and its dependencies only to build the project, uninstalling them afterwards,
and only keeping statically built binaries.
