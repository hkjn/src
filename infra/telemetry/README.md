# telemetry

This directory holds a client and server-based system to report
in various metrics from clients, like the platform they're running
on and any other context they might provide.

## Regenerating protobuf files

```
$ protoc -I report/ report/report.proto --go_out=plugins=grpc:report
```

## Building

```
$ CGO_ENABLED=0 go build -o report_client ./client/
$ CGO_ENABLED=0 go build -o report_server ./server/
```

Set `GOOS=linux GOARCH=arm` in environment to build towards armv7l.

## Running client

There's a  `report_client.service`, which can be run under systemd:

```
$ sudo cp report_client.{service,timer} /usr/lib/systemd/user/
$ sudo cp {report_client,gather_facts} /usr/local/bin/
$ systemctl --user start report_client.timer
```
