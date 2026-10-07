package slogltsv

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"sync"
)

// HandlerOptions configures an LTSVHandler.
type HandlerOptions struct {
	Level      slog.Leveler
	TimeLabel  string
	TimeFormat string
	LevelLabel string
	MsgLabel   string
}

type handlerState struct {
	mu sync.Mutex
}

type boundAttr struct {
	groups []string
	attr   slog.Attr
}

type LTSVHandler struct {
	option *HandlerOptions
	state  *handlerState
	writer io.Writer
	attrs  []boundAttr
	groups []string
}

var _ slog.Handler = (*LTSVHandler)(nil)

// NewLTSVHandler creates a handler that writes one LTSV record per line.
func NewLTSVHandler(w io.Writer, opts *HandlerOptions) *LTSVHandler {
	option := HandlerOptions{}
	if opts != nil {
		option = *opts
	}
	if option.Level == nil {
		option.Level = slog.LevelInfo
	}
	if option.TimeLabel == "" {
		option.TimeLabel = "time"
	}
	if option.TimeFormat == "" {
		option.TimeFormat = "2006-01-02T15:04:05.000000000Z0700"
	}
	if option.LevelLabel == "" {
		option.LevelLabel = "level"
	}
	if option.MsgLabel == "" {
		option.MsgLabel = "msg"
	}
	return &LTSVHandler{
		option: &option,
		state:  &handlerState{},
		writer: w,
	}
}

func (h *LTSVHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *LTSVHandler) Handle(_ context.Context, record slog.Record) error {
	fields := make([]string, 0, 3+len(h.attrs)+record.NumAttrs())
	if !record.Time.IsZero() {
		fields = append(fields, escapeKey(h.option.TimeLabel)+":"+escape(record.Time.Format(h.option.TimeFormat)))
	}
	fields = append(fields, escapeKey(h.option.LevelLabel)+":"+escape(record.Level.String()), escapeKey(h.option.MsgLabel)+":"+escape(record.Message))

	for _, attr := range h.attrs {
		appendAttr(&fields, attr.groups, attr.attr)
	}
	record.Attrs(func(attr slog.Attr) bool {
		appendAttr(&fields, h.groups, attr)
		return true
	})

	line := strings.Join(fields, "\t") + "\n"
	h.state.mu.Lock()
	defer h.state.mu.Unlock()
	n, err := io.WriteString(h.writer, line)
	if err == nil && n != len(line) {
		return io.ErrShortWrite
	}
	return err
}

func (h *LTSVHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	all := make([]boundAttr, 0, len(h.attrs)+len(attrs))
	all = append(all, h.attrs...)
	for _, attr := range attrs {
		groups := append([]string(nil), h.groups...)
		all = append(all, boundAttr{groups: groups, attr: attr})
	}
	return &LTSVHandler{option: h.option, state: h.state, writer: h.writer, attrs: all, groups: h.groups}
}

func (h *LTSVHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	groups := make([]string, 0, len(h.groups)+1)
	groups = append(groups, h.groups...)
	groups = append(groups, name)
	return &LTSVHandler{option: h.option, state: h.state, writer: h.writer, attrs: h.attrs, groups: groups}
}

func appendAttr(fields *[]string, groups []string, attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return
	}
	if attr.Value.Kind() == slog.KindGroup {
		nested := append(append([]string(nil), groups...), attr.Key)
		if attr.Key == "" {
			nested = groups
		}
		for _, child := range attr.Value.Group() {
			appendAttr(fields, nested, child)
		}
		return
	}
	if attr.Key == "" {
		return
	}
	keyParts := make([]string, 0, len(groups)+1)
	for _, group := range groups {
		if group != "" {
			keyParts = append(keyParts, group)
		}
	}
	keyParts = append(keyParts, attr.Key)
	if len(keyParts) == 0 {
		return
	}
	*fields = append(*fields, escapeKey(strings.Join(keyParts, "."))+":"+escape(attr.Value.String()))
}

func escapeKey(key string) string {
	return strings.ReplaceAll(escape(key), ":", `\:`)
}

// Escape LTSV delimiters and line breaks so a value cannot create extra fields
// or records. Backslashes are doubled to keep the representation unambiguous.
func escape(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, "\t", `\t`)
	value = strings.ReplaceAll(value, "\n", `\n`)
	value = strings.ReplaceAll(value, "\r", `\r`)
	return value
}
