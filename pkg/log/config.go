package log

import (
	"fmt"
	"log/slog"
	"strings"
)

// Format represents the logging format (JSON or Text).
type Format int

const (
	FormatJSON Format = iota + 1
	FormatText
)

// Config holds configuration for the logger.
type Config struct {
	Level     slog.Level
	Format    Format
	AddSource bool
}

// UnmarshalText implements [encoding.TextUnmarshaler].
// It unmarshals the text to a log format.
func (f *Format) UnmarshalText(text []byte) error {
	switch strings.ToUpper(string(text)) {
	case "JSON":
		*f = FormatJSON
	case "TEXT":
		*f = FormatText
	default:
		return fmt.Errorf("unknown log format: %s", text)
	}
	return nil
}
