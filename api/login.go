package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"github.com/ChileKasoka/nhl-mis/util"
)

func (server *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {

	godotenv.Load("env")

	jwtSecret := os.Getenv("JWT_SECRET")

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type LoginResponse struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Email        string `json:"email"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	loginRequest := LoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	// Verify login credentials and retrieve the user from the database
	user, err := server.store.GetUserByEmail(r.Context(), loginRequest.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	// Check if the user exists and the password is correct
	err = util.ComparePass(loginRequest.Password, user.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Set expiration time for access token
	accessTokenExpirationTime := time.Now().Add(2 * time.Minute)

	// Set expiration time for refresh token
	refreshTokenExpirationTime := time.Now().Add(60 * 24 * time.Hour)

	//set jwt claims for access token
	accessClaims := jwt.MapClaims{
		"iss": "rss-access",
		"sub": user.ID,
		"iat": jwt.NewNumericDate(time.Now().UTC()),
		"exp": jwt.NewNumericDate(accessTokenExpirationTime.UTC()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set claims for refresh token
	refreshClaims := jwt.MapClaims{
		"iss": "rss-refresh",
		"sub": user.ID,
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(refreshTokenExpirationTime),
	}

	// Create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response object
	res := LoginResponse{
		ID:           user.ID.String(),
		Name:         user.Name,
		Email:        user.Email,
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}

	// Send the login response
	RespondWithJSON(w, http.StatusOK, res)
}