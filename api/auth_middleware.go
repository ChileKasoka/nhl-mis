package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type contextKey string

const (
	accessTokenKey contextKey = "accessToken"
	claimsKey      contextKey = "claims"
)

func (server *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			godotenv.Load("env")

			// Get the JWT secret from the environment variable
			jwtSecret := os.Getenv("JWT_SECRET")
			return []byte(jwtSecret), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized, Not Valid", http.StatusUnauthorized)
			return
		}

		// Store the token and claims in the request context
		ctx := context.WithValue(r.Context(), accessTokenKey, tokenString)
		ctx = context.WithValue(ctx, claimsKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}