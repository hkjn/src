# dashboard

Package dashboard implements a web dashboard for monitoring.

See docs at http://hkjn.me/dashboard.

## Development

You can build a new binary and run it in debug mode (no auth and email
notifications) with:

```
$ go build cmd/gomon/gomon.go
$ DASHBOARD_DEBUG=true ./gomon
```

## TODO

- add icmp / ping probe package to check reachability, or
  plain tcp probe, so that network partitioning / unusual
  routing issues causing indefinite delays, high packet
  loss etc can be detected
