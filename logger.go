package main

import (
	"log/slog"
	"log"
	"os"
	"io"
	"time"
)

func modifyAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		a.Value = slog.StringValue(a.Value.Time().Format(time.DateTime))
	}
	return a
}

func writter_loging(path string) io.Writer {
	file, err := os.OpenFile(path, os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Can't create log file")
	}
	w := io.MultiWriter(file, os.Stderr)
	return (w)
}

func set_logger() {
	w := writter_loging("tmp/nomen.log")
	opts := slog.HandlerOptions {
		Level: slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: modifyAttr,
	}
	handler := slog.NewJSONHandler(w, &opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
