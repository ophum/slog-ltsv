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
	handler := slogltsv.NewLTSVHandler(os.Stdout, slogltsv.Option{
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

`NewLTSVHandler` takes an `io.Writer` as its output destination. You can optionally set the minimum log level with `Option`. The default level is `slog.LevelInfo`.

```go
handler := slogltsv.NewLTSVHandler(os.Stdout, slogltsv.Option{
	Level: slog.LevelWarn,
})
```

Group attributes are written using the `group.key` format. Tabs, line breaks, carriage returns, and backslashes in values are escaped. Colons in keys are also escaped.

## License

MIT License. See [LICENSE](LICENSE) for details.
