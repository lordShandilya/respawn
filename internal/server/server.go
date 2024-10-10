package server

import (
	"context"
	"log"
	"net/http"

	"backend/internal/config"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	router := chi.NewRouter()

	//Middlewares

	setRoutes(router)

	s := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	return &Server{HttpServer: s}
}

func setRoutes(r *chi.Mux) {
	r.Route("/", func(app chi.Router) {
		app.Post("/signup")
	})
}

func (s *Server) Start() error {
	log.Printf("The server is starting on %s", s.HttpServer.Addr)
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Printf("Shutting Down Server.....")
	return s.HttpServer.Shutdown(ctx)
}
