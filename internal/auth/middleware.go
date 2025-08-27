package auth

import (
	"context"
	"net/http"

	"github.com/eduzgun/gke-microservices/internal/session"
)

func AuthMiddleware(sessionClient *session.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Unauthorized: No session found", http.StatusUnauthorized)
				return
			}

			valid, userID, err := sessionClient.ValidateSession(r.Context(), cookie.Value)
			if err != nil || !valid {
				http.SetCookie(w, &http.Cookie{
					Name:   "session_id",
					Value:  "",
					MaxAge: -1,
					Path:   "/",
				})

				http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", userID)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
