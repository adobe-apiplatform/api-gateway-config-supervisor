# api-gateway-config-supervisor
Syncs config files from Amazon S3, Google Drive, Dropbox, Amazon Cloud Drive reloading the gateway with the updates

Table of Contents
=================

* [Status](#status)
* [Quick Start](#quick-start)
* [Dependencies](#dependencies)
* [Developer Guide](#developer-guide)

Status
======
This module is in design phase at the moment.

Quick Start
============
This module should be executed alongside the gateway in order to keep track of the configuration files and reload the gateway when there is a change.

It should expose a simple REST API for the API Gateway to check during health-checks so that in case this modules fails, the gateway may also appear unhealthy. 

Dependencies
============
To be reviewed:
* https://gowalker.org/github.com/ncw/rclone
* TBD

Developer guide
===============
TBD
