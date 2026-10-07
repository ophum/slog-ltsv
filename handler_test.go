package slogltsv

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestHandleIgnoresEmptyKeyAttributes(t *testing.T) {
	tests := []struct {
		name    string
		group   string
		wantLog string
	}{
		{
			name:    "without group",
			wantLog: "level:INFO\tmsg:hello\n",
		},
		{
			name:    "within group",
			group:   "g",
			wantLog: "level:INFO\tmsg:hello\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var output bytes.Buffer
			handler := NewLTSVHandler(&output, &HandlerOptions{})
			if tt.group != "" {
				handler = handler.WithGroup(tt.group).(*LTSVHandler)
			}

			record := slog.NewRecord(time.Time{}, slog.LevelInfo, "hello", 0)
			record.AddAttrs(slog.String("", "value"))
			if err := handler.Handle(context.Background(), record); err != nil {
				t.Fatalf("Handle() error = %v", err)
			}

			if got := output.String(); got != tt.wantLog {
				t.Errorf("log output = %q, want %q", got, tt.wantLog)
			}
		})
	}
}

func TestHandleCustomLabels(t *testing.T) {
	var output bytes.Buffer
	handler := NewLTSVHandler(&output, &HandlerOptions{
		TimeLabel:  "timestamp",
		TimeFormat: time.RFC3339,
		LevelLabel: "severity",
		MsgLabel:   "message",
	})
	record := slog.NewRecord(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC), slog.LevelInfo, "hello", 0)
	if err := handler.Handle(context.Background(), record); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	const want = "timestamp:2024-01-02T03:04:05Z\tseverity:INFO\tmessage:hello\n"
	if got := output.String(); got != want {
		t.Errorf("log output = %q, want %q", got, want)
	}
}
