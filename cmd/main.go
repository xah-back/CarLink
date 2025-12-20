package main

import (
	"log/slog"
	"os"

	"github.com/mutsaevz/team-5-ambitious/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	db := config.SetUpDatabaseConnection(logger)

	if db == nil {
		logger.Error("DB is nil")
		return
	}

	logger.Info("Application started successfully")
}
