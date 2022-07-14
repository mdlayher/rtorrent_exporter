rtorrent_exporter [![GoDoc](http://godoc.org/github.com/mdlayher/rtorrent_exporter?status.svg)](http://godoc.org/github.com/mdlayher/rtorrent_exporter) [![Build Status](https://travis-ci.org/mdlayher/rtorrent_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/rtorrent_exporter) [![Coverage Status](https://coveralls.io/repos/mdlayher/rtorrent_exporter/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/rtorrent_exporter?branch=master)
=================

Command `rtorrent_exporter` provides a Prometheus exporter for rTorrent.

Package `rtorrentexporter` provides the Exporter type used in the `rtorrent_exporter`
Prometheus exporter.

MIT Licensed.

Usage
-----

Available flags for `rtorrent_exporter` include:

```
$ ./rtorrent_exporter -h
Usage of ./rtorrent_exporter:
  -rtorrent.addr string
        address of rTorrent XML-RPC server
  - rtorrent.username string
        [optional] username used for HTTP Basic authentication with rTorrent XML-RPC server
  - rtorrent.password string
        [optional] password used for HTTP Basic authentication with rTorrent XML-RPC server
  -telemetry.addr string
        host:port for rTorrent exporter (default ":9135")
  -telemetry.path string
        URL path for surfacing collected metrics (default "/metrics")
```

An example of using `rtorrent_exporter`:

```
$ ./rtorrent_exporter -rtorrent.addr http://127.0.0.1/RPC2
2016/03/09 17:39:40 starting rTorrent exporter on ":9135" for server "http://127.0.0.1/RPC2"
```

podman run --rm -it -p 9135:9135 --name rtorrent_exporter rtorrent_exporter -rtorrent.addr http://192.168.1.193:8282/RPC2 -rtorrent.username <username> -rtorrent.password <password>


Docker
------

```
docker build -t rtorrent_exporter .
docker run --rm -d -p 9135:9135 rtorrent_exporter -rtorrent.addr "http://127.0.0.1/RPC2"
```

Sample
------

Here is a screenshot of a sample dashboard created using [`grafana`](https://github.com/grafana/grafana)
with metrics from exported from `rtorrent_exporter`.

![sample](https://cloud.githubusercontent.com/assets/1926905/13891308/bad263be-ed26-11e5-9601-9d770d95c538.png)
