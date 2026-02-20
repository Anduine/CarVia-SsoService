package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

func AuthMiddleware(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Debug("Заголовок авторизації відсутній", "Authorization", authHeader)
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := ParseToken(tokenString)
		if err != nil {
			slog.Debug("Помилка авторизації", "err", err.Error())
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", token.UserID)
		ctx = context.WithValue(ctx, "username", token.Username)
		next(w, r.WithContext(ctx))
	})
}

func AuthMiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.Debug("Заголовок авторизації відсутній", "Authorization", authHeader)
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := ParseToken(tokenString)
		if err != nil {
			slog.Debug("Помилка авторизації", "err", err.Error())
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", token.UserID)
		ctx = context.WithValue(ctx, "username", token.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
