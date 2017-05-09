# Greedy
Fast and Modular Collectd HTTP Gateway

> Greedee swims over the internets to eat small metrics

![greedy](greedee.png)

[![Build Status](https://travis-ci.org/yanc0/greedee.svg?branch=master)](https://travis-ci.org/yanc0/greedee)

## How it work?

Greedee is a fast HTTP Gateway for collectd's 
`write_http` plugin. It is very modular and can be very easily
enhanced.

It supports:

* Basic Auth
* Graphite Backend
* Console Backend (debug)

## Install

* Download latest release on [release page]("https://github.com/yanc0/collectd-http-server/releases")
* Move it on `/usr/bin/greedee`
* Configure it on `/etc/greedee/config.toml` [example here](config.toml)
* Run it by executing the greedee binary
* Install collectd (tested on 5.7.1 but should work with not 
too old previous version).
* Load and configure plugin `write_http` in JSON
Output [Official docs here](https://collectd.org/wiki/index.php/Plugin:Write_HTTP#JSON_Example)

## Build

Ensure your 1.8+ Golang environment is properly setup

```
go get -u github.com/yanc0/greedee
cd $GOPATH/src/github.com/yanc0/greedee
make setup
make build
```
## Changelog

### v0.1.0 - 2017-05-08

* Plugins
  * Graphite
  * Console
* Basic Auth
* TOML Config file
* Asynchronous metrics send
* Proper logging

## To Do

- [ ] Tests
- [ ] More plugins

## Contibutors

Feel free to make a pull request

## Licence

```
The MIT License (MIT)

Copyright (c) 2016 Yann Coleu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

