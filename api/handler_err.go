package api

import (
	"net/http"
)

func (server *Server) HandlerErr(w http.ResponseWriter, r *http.Request) {
	RespondWithError(w, 400, "something went wrong")
}