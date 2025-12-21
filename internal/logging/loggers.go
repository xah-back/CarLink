package logging

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	_ = os.MkdirAll("tmp", 0755)

	file, err := os.OpenFile(
		"tmp/app.log",
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(handler)
}
