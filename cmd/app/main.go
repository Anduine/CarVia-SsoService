package main

import (
	"log/slog"
	"os"
	"sso/internal/app"
	"sso/internal/config"
	pluslog "sso/internal/lib/logger"
)

func main() {

	cfg := config.MustLoadConfig()

	log := setupPlusSlog()

	app.Run(log, cfg.Port, cfg.DBConnector, cfg.Timeout, cfg.TokenTTL)
}

func setupPlusSlog() *slog.Logger {
	opts := pluslog.PlusHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPlusHandler(os.Stdout)

	return slog.New(handler)
}
