package routes

import (
	"net/http"
	"os"
	"strings"

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
	authMW := auth.AuthMiddleware(sessionClient)

	r.mux.Handle("GET /philosophers", authMW(http.HandlerFunc(pc.HandleGetPhilosophers)))
	r.mux.Handle("GET /philosophers/", authMW(http.HandlerFunc(pc.HandleGetPhilosophers)))
	r.mux.Handle("GET /philosophers/{id}", authMW(http.HandlerFunc(pc.HandleGetPhilosopher)))
	r.mux.Handle("POST /philosophers/add", authMW(http.HandlerFunc(pc.HandleCreatePhilosopher)))

	// Interaction routes
	r.mux.Handle("POST /philosophers/{id}/comments", authMW(http.HandlerFunc(ic.HandleCreateComment)))
	r.mux.Handle("POST /philosophers/{id}/like", authMW(http.HandlerFunc(ic.HandleToggleLike)))
	r.mux.Handle("GET /philosophers/{id}/interactions", authMW(http.HandlerFunc(ic.HandleGetInteractions)))
}

func (r *Router) RegisterAuthRoutes(ac auth.AuthController, sessionClient *session.Client) {
	// Public routes
	r.mux.HandleFunc("POST /auth/login", ac.Login)
	r.mux.HandleFunc("POST /auth/register", ac.Register)

	authMW := auth.AuthMiddleware(sessionClient)

	r.mux.Handle("GET /auth/profile", authMW(http.HandlerFunc(ac.Profile)))
	r.mux.Handle("POST /auth/logout", authMW(http.HandlerFunc(ac.Logout)))
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply CORS middleware to all requests
	corsMiddleware(r.mux).ServeHTTP(w, req)
}

func getCORSOrigins() []string {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		// Fallback for local dev
		return []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}
	return strings.Split(origins, ",")
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
