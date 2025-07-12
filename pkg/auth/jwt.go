package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type JWTToken struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	jwt.StandardClaims
}

func CreateToken(username string, userID int, tokenTTL time.Duration) (string, error) {
	expirationTime := time.Now().Add(tokenTTL).Add(time.Hour).Unix()

	claims := &JWTToken{
		Username: username,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			Issuer:    "sso_service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	tk, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("could not create token: %v", err)
	}

	return tk, nil
}

func ParseToken(tokenStr string) (*JWTToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTToken)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
