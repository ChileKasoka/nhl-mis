package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	sqlc "github.com/ChileKasoka/nhl-mis/db/sqlc"
	"github.com/ChileKasoka/nhl-mis/utils"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (server *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type Users struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userparams := Users{}
	err := decoder.Decode(&userparams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "could not decode")
		return
	}

	hashedPassword, err := utils.HashedPass(userparams.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't hash password")
		return
	}
 
	//uuidValue := uuid.New()
	user, err := server.store.CreateUser(r.Context(), sqlc.CreateUserParams{
		ID:        uuid.New(),
		Name:      userparams.Name,
		Email:     userparams.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})

	res := UserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "user failed to create")
		return
	}

	RespondWithJSON(w, http.StatusOK, res)
}

func (server *Server) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the user ID from the request path using chi
	id := chi.URLParam(r, "id")

	// Parse the ID string into a UUID
	userID, err := uuid.Parse(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Call the getUser function to retrieve the user from the database
	user, err := server.store.GetUser(r.Context(), sqlc.GetUserParams{ID: userID})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to retrieve user")
		return
	}

	// Create a response object
	res := UserResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}

	// Send the user response
	RespondWithJSON(w, http.StatusOK, res)
}

// Update User handler
func (server *Server) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the authorization token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		RespondWithError(w, http.StatusUnauthorized, "missing authorization token")
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate the authorization token
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Return the secret key used to sign the token
		jwtSecret := os.Getenv("JWT_SECRET")
		return []byte(jwtSecret), nil
	})
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid authorization token")
		return
	}

	// Verify if the token is valid and contains the required claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		RespondWithError(w, http.StatusUnauthorized, "failed to verify")
		return
	}

	// Extract the user ID from the claims
	id, ok := claims["sub"].(string)
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "invalid user id")
		return
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	// Get the user from the database using the id
	user, err := server.store.GetUser(r.Context(), sqlc.GetUserParams{ID: userID})
	if err != nil {
		// Handle the error, e.g., log it, return an error response, etc.
		log.Printf("Error getting user: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}
	// Implement logic to update the user fields like name and email

	type UpdateUserReq struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password,omitempty"`
	}

	var updateReq UpdateUserReq
	err = json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "failed to decode request body")
		return
	}

	// Update the user fields only if they are not empty
	if updateReq.Name != "" {
		user.Name = updateReq.Name
	}
	if updateReq.Email != "" {
		user.Email = updateReq.Email
	}

	// Only update the password if it's not empty
	if updateReq.Password != "" {
		hashedPassword, err := utils.HashedPass(updateReq.Password)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "couldn't hash password")
			return
		}
		user.Password = hashedPassword
	}

	// Create the UpdateUserParams with the updated values
	updateParams := sqlc.UpdateUserParams{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	if err != nil {
		// Handle the error, e.g., log it, return an error response, etc.
		log.Printf("Error creating updateParams: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to create updateParams")
		return
	}

	updatedUser, err := server.store.UpdateUser(r.Context(), updateParams)

	if err != nil {
		// Handle the error, e.g., log it, return an error response, etc.
		log.Printf("Error creating updateParams: %v", err)
		RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Return a success response
	RespondWithJSON(w, http.StatusOK, updatedUser)
}