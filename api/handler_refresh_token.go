package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func (server *Server) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the refresh token from the request
	type RefreshResponse struct {
		Token string `json:"token"`
	}

	refreshToken := r.Header.Get("Authorization")
	if refreshToken == "" {
		RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	// Extract the token string from the Authorization header
	tokenString := strings.TrimPrefix(refreshToken, "Bearer ")

	// Verify the refresh token and extract the user ID
	userID, err := server.ValidateRefreshToken(tokenString)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// Generate a new access token
	accessToken, err := server.GenerateAccessToken(userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to generate access token")
		return
	}

	// Create the JWT claims with the user ID
	// claims := jwt.MapClaims{
	// 	"sub": userID,
	// }

	// Return the new access token to the client
	// RespondWithJSON(w, http.StatusOK, map[string]string{
	// 	"access_token": accessToken,
	// })

	// Create the refresh response
	res := RefreshResponse{
		Token: accessToken,
	}

	// Return the new access token
	json.NewEncoder(w).Encode(res)

	w.Header().Set("Authorization", "Bearer "+accessToken)

	// Call the next handler in the chain
	// 	ctx := context.WithValue(r.Context(), "accessToken", accessToken)
	// 	ctx = context.WithValue(ctx, "claims", claims)
	// 	next.ServeHTTP(w, r.WithContext(ctx))
	// })
}

// ValidateRefreshToken validates the refresh token and returns the associated user ID if it is valid
func (server *Server) ValidateRefreshToken(refreshToken string) (string, error) {
	// Parse and validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid refresh token signing method")
		}

		godotenv.Load("env")

		// Get the JWT secret from the environment variable
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", errors.New("could not parse and validate")
	}

	// Verify if the token is valid and contains the required claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	// Extract the user ID from the claims
	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid user")
	}

	if claims["iss"] != "rss-refresh" {
		return "", errors.New("needs to be refresh token")
	}

	return userID, nil
}

// GenerateAccessToken generates a new access token for the provided user ID
func (server *Server) GenerateAccessToken(userID string) (string, error) {
	// Create the JWT claims with the user ID
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(), // Set the expiration time for the access token
		// Add any other claims you want to include in the access token
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get the JWT secret from the environment variable
	jwtSecret := os.Getenv("JWT_SECRET")

	// Sign the token with the JWT secret
	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}