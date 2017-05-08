# collectd-http-server
Fast and Modular Collectd HTTP Gateway

## How it work?

Collectd HTTP Server is a fast HTTP Gateway for collectd's 
`write_http` plugin. It is very modular and can be very easily
enhanced.

It supports:

* Basic Auth
* Graphite Backend
* Console Backend (debug)

## Install

* Download latest release on [release page]("https://github.com/yanc0/collectd-http-server/releases")
* Move it on `/usr/bin/collectd-http-server`
* Configure it on `/etc/collectd-http-server/config.toml` [example here](config.toml)
* Run it by executing the collectd-http-server binary
* Install collectd (tested on 5.7.1 but should work with not 
too old previous version).
* Load and configure plugin `write_http` in JSON
Output [Official docs here](https://collectd.org/wiki/index.php/Plugin:Write_HTTP#JSON_Example)

## Build

Ensure your 1.8+ Golang environment is properly setup

```
go get -u github.com/yanc0/collectd-http-server
cd $GOPATH/src/github.com/yanc0/collectd-http-server
make setup
make build
```

## Contibutors
