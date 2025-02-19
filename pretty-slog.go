package prettylogger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type LogHandlerOptions struct {
	SlogOptions *slog.HandlerOptions
}

type LogHandler struct {
	opts  LogHandlerOptions
	inner slog.Handler
	l     *slog.Logger
	attrs []slog.Attr
}

func (opts LogHandlerOptions) NewLogHandler(out io.Writer) *LogHandler {
	inner := slog.NewJSONHandler(out, opts.SlogOptions)
	return &LogHandler{
		opts:  opts,
		inner: inner,
		l:     slog.New(inner),
	}
}

func (h *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := r.Time.Format("[15:04:05]")
	msg := color.CyanString(r.Message)

	fmt.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		opts:  h.opts,
		inner: h.inner,
		l:     h.l,
		attrs: append(h.attrs, attrs...),
	}
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		opts:  h.opts,
		inner: h.inner.WithGroup(name),
		l:     h.l,
		attrs: h.attrs,
	}
}
