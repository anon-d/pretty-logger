package prettylogger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type LogHandlerOptions struct {
	SlogOptions *slog.HandlerOptions
}

type LogHandler struct {
	opts LogHandlerOptions
	slog.Handler
	l     *log.Logger
	attrs []slog.Attr
}

func (opts LogHandlerOptions) NewLogHandler(out io.Writer) *LogHandler {
	return &LogHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOptions),
		l:       log.New(out, "", 0),
	}
}

func (h *LogHandler) Handle(_ context.Context, r slog.Record) error {
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

	timeStr := r.Time.Format("[15:05:05]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

func (h *LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogHandler{
		Handler: h.Handler,
		l:       h.l,
		attrs:   attrs,
	}
}

func (h *LogHandler) WithGroup(name string) slog.Handler {
	return &LogHandler{
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
	}
}
