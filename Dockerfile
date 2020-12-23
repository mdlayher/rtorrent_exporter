# Builder image
FROM docker.io/golang:1.13 AS builder

ENV GO111MODULE=on

EXPOSE 9135

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...


# Final image
FROM docker.io/debian:buster-slim

RUN groupadd -g 911 -r app && useradd -u 911 --no-log-init -r -g app app
COPY --chown=app:app --from=builder /go/bin/rtorrent_exporter .

ENTRYPOINT ["/rtorrent_exporter"]
