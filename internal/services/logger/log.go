// internal/logger/logger.go
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

var Log *slog.Logger

func Init(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	Log = slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Any().(time.Time)
				a.Value = slog.StringValue(t.Format("2006/01/02 15:04:05"))
			}
			return a
		},
	}))

	return nil
}

// helper biar kayak Laravel
func Error(msg string, args ...any) {
	Log.Error(fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any) {
	Log.Info(fmt.Sprintf(msg, args...))
}

func Debug(msg string, args ...any) {
	Log.Debug(fmt.Sprintf(msg, args...))
}
