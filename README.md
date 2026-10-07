# slog-ltsv

`slog-ltsv` は、Go 標準の [`log/slog`](https://pkg.go.dev/log/slog) 用 LTSV ハンドラーです。ログを 1 行 1 レコードで出力し、属性をタブ区切りの `key:value` 形式にします。

## インストール

```sh
go get github.com/ophum/slog-ltsv
```

## 使い方

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

出力例:

```text
time:2026-10-07T12:34:56.123456789Z	level:INFO	msg:server started	service:checkout	addr::8080
time:2026-10-07T12:34:56.123456790Z	level:INFO	msg:request received	service:checkout	request.id:req-123	request.method:GET
```

実際の時刻やタイムゾーンは実行環境によって異なります。

## 設定

`NewLTSVHandler` は出力先の `io.Writer` を受け取ります。省略可能な `Option` で最低ログレベルを指定できます。既定値は `slog.LevelInfo` です。

```go
handler := slogltsv.NewLTSVHandler(os.Stdout, slogltsv.Option{
	Level: slog.LevelWarn,
})
```

グループ属性は `group.key` の形式で出力します。値に含まれるタブ、改行、復帰、バックスラッシュはエスケープされ、キーに含まれるコロンもエスケープされます。

## ライセンス

MIT License. 詳細は [LICENSE](LICENSE) を参照してください。
