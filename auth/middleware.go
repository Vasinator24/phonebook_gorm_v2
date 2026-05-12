package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey = contextKey("user")
const secret = "my_secret_key"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*Claims)

		ctx := context.WithValue(r.Context(), UserKey, claims)
		next(w, r.WithContext(ctx))
	}
}

func GetUserFromContext(r *http.Request) *Claims {
	claims, ok := r.Context().Value(UserKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenStr := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)

		token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		claims := token.Claims.(jwt.MapClaims)

		role := claims["role"].(string)

		if role != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
