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
  -telemetry.addr string
        host:port for rTorrent exporter (default ":9135")
  -telemetry.path string
        URL path for surfacing collected metrics (default "/metrics")
```

Examples of using `rtorrent_exporter`:

```
$ ./rtorrent_exporter -rtorrent.addr http://127.0.0.1/RPC2
2016/03/09 17:39:40 starting rTorrent exporter on ":9135" for server "http://127.0.0.1/RPC2"
```

```
$ ./rtorrent_exporter -rtorrent.addr scgi://127.0.0.1:5000
2016/03/09 17:39:40 starting rTorrent exporter on ":9135" for server "scgi://127.0.0.1:5000"
```

```
$ ./rtorrent_exporter -rtorrent.addr scgi:///home/user/rtorrent/rpc.socket
2016/03/09 17:39:40 starting rTorrent exporter on ":9135" for server "scgi:///home/user/rtorrent/rpc.socket"
```

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
