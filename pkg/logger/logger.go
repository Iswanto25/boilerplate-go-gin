package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	fileLog     *os.File
	fileLogMu   sync.Mutex
	currentDate string
	logDir      string
)

func InitLogger(env string) {
	logDir = filepath.Join(getCWD(), "logger")
	rotateFile()

	var consoleHandler slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}

	if env == "production" {
		consoleHandler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		consoleHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.Attr{}
				}
				if a.Key == slog.LevelKey {
					return slog.Attr{}
				}
				if a.Key == slog.SourceKey {
					return slog.Attr{}
				}
				return a
			},
		})
	}

	fileHandler := slog.NewJSONHandler(newDailyFileWriter(), &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(teeHandler{console: consoleHandler, file: fileHandler})
	slog.SetDefault(logger)
}

// LogConsole writes a simplified log line in Express.js style:
// {METHOD} {PATH} {STATUS} | {userName} | {responseTime}ms
func LogConsole(method, path string, status int, userName string, responseTimeMs int64) {
	level := "INFO"
	if status >= 400 {
		level = "ERROR"
	}
	line := fmt.Sprintf("[%s] %s %s %d | %s | %dms", level, method, path, status, userName, responseTimeMs)
	slog.Info(line)
}

type teeHandler struct {
	console slog.Handler
	file    slog.Handler
}

func (t teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return t.console.Enabled(ctx, level) || t.file.Enabled(ctx, level)
}

func (t teeHandler) Handle(ctx context.Context, r slog.Record) error {
	if t.console.Enabled(ctx, r.Level) {
		_ = t.console.Handle(ctx, r.Clone())
	}
	if t.file.Enabled(ctx, r.Level) {
		_ = t.file.Handle(ctx, r.Clone())
	}
	return nil
}

func (t teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return teeHandler{
		console: t.console.WithAttrs(attrs),
		file:    t.file.WithAttrs(attrs),
	}
}

func (t teeHandler) WithGroup(name string) slog.Handler {
	return teeHandler{
		console: t.console.WithGroup(name),
		file:    t.file.WithGroup(name),
	}
}

type dailyFileWriter struct{}

func newDailyFileWriter() *dailyFileWriter {
	return &dailyFileWriter{}
}

func (w *dailyFileWriter) Write(p []byte) (n int, err error) {
	fileLogMu.Lock()
	defer fileLogMu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today != currentDate || fileLog == nil {
		if fileLog != nil {
			fileLog.Close()
		}
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fileLog = nil
			return len(p), nil
		}
		filePath := filepath.Join(logDir, today+".log")
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fileLog = nil
			return len(p), nil
		}
		fileLog = f
		currentDate = today
	}

	if fileLog != nil {
		return fileLog.Write(p)
	}
	return len(p), nil
}

func rotateFile() {
	fileLogMu.Lock()
	defer fileLogMu.Unlock()

	today := time.Now().Format("2006-01-02")
	if today == currentDate && fileLog != nil {
		return
	}

	if fileLog != nil {
		fileLog.Close()
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return
	}

	filePath := filepath.Join(logDir, today+".log")
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	fileLog = f
	currentDate = today
}

func getCWD() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
