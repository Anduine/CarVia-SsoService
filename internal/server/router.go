package server

import (
	"log/slog"
	"net/http"
	"sso-service/internal/delivery/http_handlers"
	"sso-service/pkg/auth"

	"github.com/gorilla/mux"
)

func NewRouter(usersHandler *http_handlers.UsersHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/sso/register", usersHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/sso/login", usersHandler.LoginHandler).Methods("POST")

	router.Handle("/api/sso/user_profile", auth.AuthMiddleware(usersHandler.UserProfileHandler)).Methods("GET")
	router.Handle("/api/sso/update_user_profile", auth.AuthMiddleware(usersHandler.UpdateUserProfileHandler)).Methods("PUT")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Маршрут не знайдено", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Маршрут не знайдено", http.StatusNotFound)
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Заборонений метод", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Заборонений метод", http.StatusMethodNotAllowed)
	})

	return router
}
