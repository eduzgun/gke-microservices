package routes

import (
	"net/http"
	"os"

	"github.com/eduzgun/gke-microservices/internal/auth"
	"github.com/eduzgun/gke-microservices/internal/interaction"
	"github.com/eduzgun/gke-microservices/internal/philosopher"
	"github.com/eduzgun/gke-microservices/internal/session"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &Router{mux: mux}
}

func (r *Router) RegisterPhilosopherRoutes(
	pc philosopher.PhilosopherController,
	ic interaction.InteractionController,
	sessionClient *session.Client,
) {
	protected := http.NewServeMux()

	// Philosopher routes
	protected.HandleFunc("GET /philosophers/", pc.HandleGetPhilosophers)
	protected.HandleFunc("GET /philosophers/{id}", pc.HandleGetPhilosopher)
	protected.HandleFunc("POST /philosophers/add", pc.HandleCreatePhilosopher)

	// Interaction routes
	protected.HandleFunc("POST /philosophers/{id}/comments", ic.HandleCreateComment)
	protected.HandleFunc("POST /philosophers/{id}/like", ic.HandleToggleLike)
	protected.HandleFunc("GET /philosophers/{id}/interactions", ic.HandleGetInteractions)

	r.mux.Handle("/philosophers", auth.AuthMiddleware(sessionClient)(protected))
	r.mux.Handle("/philosophers/", auth.AuthMiddleware(sessionClient)(protected))
}

func (r *Router) RegisterAuthRoutes(ac auth.AuthController, sessionClient *session.Client) {
	// Public routes
	r.mux.HandleFunc("POST /auth/login", ac.Login)
	r.mux.HandleFunc("POST /auth/register", ac.Register)

	protected := auth.AuthMiddleware(sessionClient)

	r.mux.Handle("GET /auth/profile", protected(http.HandlerFunc(ac.Profile)))
	r.mux.Handle("POST /auth/logout", protected(http.HandlerFunc(ac.Logout)))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply CORS middleware to all requests
	corsMiddleware(r.mux).ServeHTTP(w, req)
}

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
