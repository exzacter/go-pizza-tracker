package main

import (
	"go-pizza-tracker/internal/models"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// loaded the config
	cfg := loadConfig()

	//slog is a modern way to log in application
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dbModel, err := models.InitDB(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialise database", "error", err)
		os.Exit(1)
	}

	slog.Info("Database initialised successfully")

	RegisterCustomValidators()

	h := NewHandler(dbModel)

	router := gin.Default()

	if err := loadTemplates(router); err != nil {
		slog.Error("Failed to load templates", "error", err)
		os.Exit(1)
	}

	setupRotues(router, h)

	slog.Info("Server starting", "url", "http://localhost:"+cfg.Port)

	router.Run(":" + cfg.Port)
}
