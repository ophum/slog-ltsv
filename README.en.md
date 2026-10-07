# slog-ltsv

`slog-ltsv` is an LTSV handler for Go's standard [`log/slog`](https://pkg.go.dev/log/slog) package. It writes one log record per line, with attributes formatted as tab-separated `key:value` pairs.

## Installation

```sh
go get github.com/ophum/slog-ltsv
```

## Usage

```go
package main

import (
	"log/slog"
	"os"

	slogltsv "github.com/ophum/slog-ltsv"
)

func main() {
	handler := slogltsv.NewLTSVHandler(os.Stdout, &slogltsv.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler).With(
		slog.String("service", "checkout"),
	)

	logger.Info("server started", slog.String("addr", ":8080"))
	logger.WithGroup("request").Info("request received",
		slog.String("id", "req-123"),
		slog.String("method", "GET"),
	)
}
```

Example output:

```text
time:2026-10-07T12:34:56.123456789Z	level:INFO	msg:server started	service:checkout	addr::8080
time:2026-10-07T12:34:56.123456790Z	level:INFO	msg:request received	service:checkout	request.id:req-123	request.method:GET
```

The actual timestamps and time zone depend on the runtime environment.

## Configuration

`NewLTSVHandler` takes an `io.Writer` and a `*HandlerOptions`. Unspecified values use their defaults, including when `nil` is passed. The default minimum level is `slog.LevelInfo`; the default labels are `time`, `level`, and `msg`; the default time format is `2006-01-02T15:04:05.000000000Z0700`. `TimeFormat` uses the same layout syntax as Go's [`time.Time.Format`](https://pkg.go.dev/time#Time.Format).

```go
handler := slogltsv.NewLTSVHandler(os.Stdout, &slogltsv.HandlerOptions{
	Level:      slog.LevelWarn,
	TimeLabel:  "timestamp",
	TimeFormat: "2006-01-02T15:04:05Z07:00",
	LevelLabel: "severity",
	MsgLabel:   "message",
})
```

Group attributes are written using the `group.key` format. Tabs, line breaks, carriage returns, and backslashes in values are escaped. Colons in keys are also escaped.

## License

MIT License. See [LICENSE](LICENSE) for details.
