package prettylogger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

// const TraceLevel = slog.Level(-8)

// Debug - 4
// Info 0
// Warn 4
// Error 8

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	color.NoColor = false

	level := r.Level.String() + ":"

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

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		color.HiBlackString(string(b)),
	)

	return nil
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

// func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
// 	return &PrettyHandler{
// 		opts:    h.opts,
// 		Handler: h.Handler,
// 		l:       h.l,
// 		attrs:   append(h.attrs, attrs...),
// 	}
// }

// func (h *PrettyHandler) WithGroup(name string) slog.Handler {
// 	return &PrettyHandler{
// 		opts:    h.opts,
// 		Handler: h.Handler.WithGroup(name),
// 		l:       h.l,
// 		attrs:   h.attrs,
// 	}
// }

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
