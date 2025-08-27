// internal/routes/router.go

package routes

import (
	"net/http"
	"os"

	"github.com/eduzgun/gke-microservices/internal/auth"
	"github.com/eduzgun/gke-microservices/internal/philosopher"
	"github.com/eduzgun/gke-microservices/internal/session"
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

// RegisterPhilosopherRoutes registers philosopher-related endpoints (all protected for now)
func (r *Router) RegisterPhilosopherRoutes(pc philosopher.PhilosopherController, sessionClient *session.Client) {
	// Wrap all philosopher routes with auth middleware
	protected := http.NewServeMux()
	protected.HandleFunc("GET /philosophers", pc.HandleGetPhilosophers)
	protected.HandleFunc("GET /philosophers/{id}", pc.HandleGetPhilosopher)
	protected.HandleFunc("POST /philosophers", pc.HandleCreatePhilosopher)

	// Apply AuthMiddleware to all /philosophers routes
	r.mux.Handle("/philosophers", auth.AuthMiddleware(sessionClient)(protected))
	r.mux.Handle("/philosophers/", auth.AuthMiddleware(sessionClient)(protected))
}

// RegisterAuthRoutes registers auth endpoints (public + protected)
func (r *Router) RegisterAuthRoutes(ac auth.AuthController, sessionClient *session.Client) {
	// Public routes
	r.mux.HandleFunc("POST /auth/login", ac.Login)
	r.mux.HandleFunc("POST /auth/register", ac.Register)
	r.mux.HandleFunc("POST /auth/logout", ac.Logout)

	// Protected routes
	protected := http.NewServeMux()
	protected.HandleFunc("GET /profile", ac.Profile)

	// Apply AuthMiddleware to protected routes
	r.mux.Handle("/profile", auth.AuthMiddleware(sessionClient)(protected))
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
			"https://yourproductionapp.com",
		}
	default:
		return []string{
			"http://localhost:3000",
			"http://localhost:5173", // Vite dev
			"http://localhost:4173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}
}

// corsMiddleware adds CORS headers
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
