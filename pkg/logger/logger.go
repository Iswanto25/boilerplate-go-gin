package logger

import (
	"log/slog"
	"os"
)

// InitLogger inisialisasi default logger dari package log/slog bawaan Go 1.21+
func InitLogger(env string) {
	var handler slog.Handler

	if env == "production" {
		// Output JSON untuk production
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// Output text biasa untuk development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
