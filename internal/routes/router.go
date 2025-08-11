package routes

import (
	"net/http"
	"os"

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
	// Apply CORS middleware to all requests
	corsMiddleware(r.mux).ServeHTTP(w, req)
}

// getCORSOrigins returns allowed origins based on environment
func getCORSOrigins() []string {
	env := os.Getenv("ENVIRONMENT")

	switch env {
	case "prod":
		return []string{
			"https://deployedapp",
		}
	default:
		return []string{
			"http://localhost:3000",
			"http://localhost:5173", // Vite
			"http://localhost:4173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	allowedOrigins := getCORSOrigins()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
