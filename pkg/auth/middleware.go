package auth

import (
	"context"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Заголовок авторизації відсутній")
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			log.Println("Bearer токен відсутній")
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		claims, err := ParseToken(tokenStr)
		if err != nil {
			log.Println("Неправильний токен авторизації: ", err)
			http.Error(w, "Не авторизовано", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "login", claims.Username)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
