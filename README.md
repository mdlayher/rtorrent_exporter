rtorrent_exporter [![GoDoc](http://godoc.org/github.com/mdlayher/rtorrent_exporter?status.svg)](http://godoc.org/github.com/mdlayher/rtorrent_exporter) [![Build Status](https://travis-ci.org/mdlayher/rtorrent_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/rtorrent_exporter) [![Coverage Status](https://coveralls.io/repos/mdlayher/rtorrent_exporter/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/rtorrent_exporter?branch=master)
=================

Command `rtorrent_exporter` provides a Prometheus exporter for rTorrent.

Package `rtorrentexporter` provides the Exporter type used in the `rtorrent_exporter`
Prometheus exporter.

MIT Licensed.

# Usage

Available flags for `rtorrent_exporter` include:

```
$ ./rtorrent_exporter -h
Usage of ./rtorrent_exporter:
  -rtorrent_addr string
        address of rTorrent XML-RPC server
  -telemetry_addr string
        host:port for rTorrent exporter (default ":9135")
  -telemetry_path string
        URL path for surfacing collected metrics (default "/metrics")
```

An example of using `rtorrent_exporter`:

```
$ ./rtorrent_exporter -rtorrent_addr http://127.0.0.1/RPC2
2016/03/09 17:39:40 starting rTorrent exporter on ":9135" for server "http://127.0.0.1/RPC2"
```

You can also use environment variables instead of flags, for example:

```
$ RTORRENT_ADDR=http://127.0.0.1/RPC2 ./rtorrent_exporter
```

# Sample

Here is a screenshot of a sample dashboard created using [`grafana`](https://github.com/grafana/grafana)
with metrics from exported from `rtorrent_exporter`.

![sample](https://cloud.githubusercontent.com/assets/1926905/13891308/bad263be-ed26-11e5-9601-9d770d95c538.png)

# Building and running a Docker container

*Replace `tehwey` with your own Docker Hub name if you want to create your own image.*

**Build Docker image and push:**

```
docker build -t tehwey/docker-rtorrent-exporter .
docker tag <image ID> tehwey/docker-rtorrent-exporter
docker push tehwey/docker-rtorrent-exporter
```

**Create Docker container:**

```
docker create --name=prom-rtorrent-exporter \
-p 9005:9135 \
-e RTORRENT_ADDR=http://localhost:8005/RPC2 \
tehwey/docker-rtorrent-exporter
```

**Run Docker container:**

```
docker run --name=prom-rtorrent-exporter \
-p 9005:9135 \
-e RTORRENT_ADDR=http://localhost:8005/RPC2 \
tehwey/docker-rtorrent-exporter
```

