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

An example of using `rtorrent_exporter`:

```
$ ./rtorrent_exporter -rtorrent.addr http://127.0.0.1/RPC2
2016/03/09 17:39:40 starting rTorrent exporter on ":9135" for server "http://127.0.0.1/RPC2"
```

Sample
------

Here is a screenshot of a sample dashboard created using [`grafana`](https://github.com/grafana/grafana)
with metrics from exported from `rtorrent_exporter` (see sample_dash.json).

![sample](https://user-images.githubusercontent.com/766820/29206006-923e7b20-7e45-11e7-84ae-1e51d0755f92.jpg)
