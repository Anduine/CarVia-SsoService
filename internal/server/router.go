package server

import (
	"net/http"
	"sso/internal/delivery/http_handlers"
	"sso/pkg/auth"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func NewRouter(usersHandler *http_handlers.UsersHandler) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/sso/register", usersHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/sso/login", usersHandler.LoginHandler).Methods("POST")

	router.Handle("/api/sso/user_profile", auth.AuthMiddleware(http.HandlerFunc(usersHandler.UserProfileHandler)))
	router.Handle("/api/sso/update_user_profile", auth.AuthMiddleware(http.HandlerFunc(usersHandler.UpdateUserProfileHandler)))

	router.HandleFunc("/api/sso/images/{filename}", http_handlers.ServeUserAvatar).Methods("GET")


	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3010"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(router)
	
	return handler
}
