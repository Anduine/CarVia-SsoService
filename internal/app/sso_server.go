package app

import (
	"log/slog"
	"sso/internal/delivery/http_handlers"
	"sso/internal/repository"
	"sso/internal/server"
	"sso/internal/service"
	"sso/pkg/database"
	"time"
)

func Run(log *slog.Logger, port, dbConn string, timeout, tokenTTL time.Duration) {
	db := database.NewPostgresConnection(dbConn)

	repo := repository.NewPostgresUserRepo(db)
	usersService := service.NewUsersService(repo)
	usersHandler := http_handlers.NewUsersHandler(log, usersService, tokenTTL)

	handler := server.NewRouter(usersHandler)

	server.StartServer(log, handler, port, timeout)
}
