package api

import (
	"log"
	"net/http"
)

func (server *Server) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	RespondWithJSON(w, 200, struct{}{})
}