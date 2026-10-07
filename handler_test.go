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
			handler := NewLTSVHandler(&output)
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
