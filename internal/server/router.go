package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"sekisho/internal/handler"
)

// NewRouter creates and configures a new chi.Router.
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	
	// Basic middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint
	r.Get("/health", handler.Health)

	return r
}
