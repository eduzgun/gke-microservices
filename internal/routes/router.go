package routes

import (
	"net/http"

	"github.com/eduzgun/gke-microservices/internal/philosopher"
)

type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new HTTP router with base routes
func NewRouter() *Router {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &Router{mux: mux}
}

func (r *Router) RegisterPhilosopherRoutes(controller philosopher.PhilosopherController) {
	r.mux.HandleFunc("GET /philosophers", controller.HandleGetPhilosophers)
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
