package api

import (
	"log"
	"net/http"

	db "github.com/ChileKasoka/nhl-mis/db/sqlc"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type Server struct {
	store  *db.Store
	router *chi.Mux
}

func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Use(MiddlewareLogger)

	v1Router.Get("/healthz", server.HandlerReadiness)
	v1Router.Get("/err", server.HandlerErr)
	v1Router.Post("/users", server.CreateUserHandler)
	v1Router.Get("/users/{id}", server.GetUserHandler)
	v1Router.Post("/login", server.LoginHandler)
	v1Router.Post("/refresh", server.RefreshTokenHandler)
	v1Router.Put("/users", server.UpdateUserHandler)

	router.Mount("/v1", v1Router)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}

	return srv.ListenAndServe()
}
