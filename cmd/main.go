package main

import (
	"github.com/mutsaevz/team-5-ambitious/internal/config"
	"github.com/mutsaevz/team-5-ambitious/internal/logging"
)

func main() {
	// инициализация логгера (tmp внутри logging)
	logger := logging.New()

	db := config.SetUpDatabaseConnection(logger)
	if db == nil {
		logger.Error("database is nil")
		return
	}

	logger.Info("application started successfully")
}
