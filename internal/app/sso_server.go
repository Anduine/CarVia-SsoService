package app

import (
	"sso-service/internal/config"
	"sso-service/internal/delivery/http_handlers"
	"sso-service/internal/repository"
	"sso-service/internal/server"
	"sso-service/internal/service"
	"sso-service/pkg/database"
)

func Run(cfg *config.Config) {
	db := database.NewPostgresConnection(cfg.DB.Host, cfg.DB.DBName, cfg.DB.User, cfg.DB.Password)

	repo := repository.NewPostgresUserRepo(db)
	usersService := service.NewUsersService(repo, cfg.StorageURL)
	usersHandler := http_handlers.NewUsersHandler(usersService, cfg.TokenTTL)

	handler := server.NewRouter(usersHandler)

	server.StartServer(handler, cfg.Port, cfg.Timeout)
}
