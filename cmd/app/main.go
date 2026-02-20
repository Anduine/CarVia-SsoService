package main

import (
	"log/slog"
	"os"
	"sso-service/internal/app"
	"sso-service/internal/config"
	"sso-service/internal/lib/logger"
)

func main() {
	config := config.MustLoadConfig()

	logger.InitGlobalLogger(os.Stdout, slog.LevelDebug)

	app.Run(config)
}
