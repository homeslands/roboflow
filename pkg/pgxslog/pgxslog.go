package pgxslog

import (
	"context"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct {
	l *slog.Logger
}

func NewTracer(l *slog.Logger) pgx.QueryTracer {
	return &tracelog.TraceLog{
		Logger:   &Logger{l: l},
		LogLevel: tracelog.LogLevelDebug,
	}
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	logger := l.l
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	// remove \n, \t from msg
	m := strings.ReplaceAll(msg, "\n", " ")
	formattedMsg := strings.ReplaceAll(m, "\t", " ")

	logger.LogAttrs(ctx, translateLevel(level), formattedMsg, attrs...)
}

func translateLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace:
		return slog.LevelDebug
	case tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	case tracelog.LogLevelError:
		return slog.LevelError
	case tracelog.LogLevelNone:
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
