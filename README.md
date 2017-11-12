# Greedee
Fast and Modular Collectd HTTP Gateway

> Greedee swims over the internets to eat small metrics and events

![greedy](greedee.png)

[![Build Status](https://travis-ci.org/yanc0/greedee.svg?branch=master)](https://travis-ci.org/yanc0/greedee)

## How it work?

Greedee is a fast HTTP Gateway for collectd's 
`write_http` plugin. It is very modular and can be very easily
enhanced. It can also ingests events (notifications) to be
stored on database

It supports:

* Basic Auth
* Graphite Backend
* MySQL Backend (for events)
* Console Backend (debug)

### Configure Collectd
```
<Plugin write_http>
   <Node "greedee">
       URL "http://127.0.0.1:9223/metrics"
       User "user"
       Password "pass"
       Format "JSON"
   </Node>
</Plugin>
```

### Sending an event through the API
```
curl -d '{"name": "postgresql_backup",
          "ttl": 36000,
          "status": 0,
          "source": "curl",
          "description": "Nightly postgresql backup"}' \
          http://127.0.0.1:9223/events -u "user:pass"
```

**name (mandatory)**: name of the event
**ttl**: time to live in seconds, if set, trigger failed event if no event with the same name is created after this time
**status**: OK = 0, FAILURE > 0
**description**: Short explanation of the event
**source**: Source of event (example: curl, cron, ansible, etc.)

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

### v0.3.2 - 2017-11-12

* Remove 00-00-0000 date for mysql 5.7 strict mode compatibility

### v0.3.1 - 2017-10-23

* Add route /version

### v0.3.0 - 2017-10-09

* Event Reactor - Trigger failed events when events are expired
* Metric Transformer - Calculate derive metrics ith previous data


### v0.2.0 - 2017-06-28

* Events support
* Plugins
  * MySQL (events)
* Some code clean

### v0.1.0 - 2017-05-08

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

## Contributors

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

