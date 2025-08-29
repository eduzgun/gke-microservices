package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eduzgun/gke-microservices/internal/errs"
	"github.com/eduzgun/gke-microservices/internal/models"
	"github.com/eduzgun/gke-microservices/internal/user"
)

type authControllerImpl struct {
	authService AuthService
	userService user.UserService
	logger      *slog.Logger
}

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Profile(w http.ResponseWriter, r *http.Request)
}

func NewAuthController(
	authService AuthService,
	userService user.UserService, // Interface so no pointer needed
	logger *slog.Logger,
) AuthController {
	return &authControllerImpl{
		authService: authService,
		userService: userService,
		logger:      logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

func (ac *authControllerImpl) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := ac.userService.CreateUser(r.Context(), user); err != nil {
		ac.logger.Error("failed to create user", "err", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(LoginResponse{
		Status:  "success",
		Message: "User registered successfully",
	})
}

func (ac *authControllerImpl) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	sessionID, err := ac.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// Check for which type of error it is
		switch {
		case errors.Is(err, errs.ErrInvalidCredentials):
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		}
		http.Error(w, "Server failed on Login, please try again later", http.StatusInternalServerError)
		ac.logger.Error("failed to login", "err", err)
		return
	}

	// Set secure cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Status:  "success",
		Message: "Login successful",
	})
}

func (ac *authControllerImpl) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Status:  "success",
			Message: "Already logged out",
		})
		return
	}

	// Invalidate session remotely
	err = ac.authService.Logout(r.Context(), cookie.Value)
	if err != nil {
		// Log error but don't fail the logout
		// Could be network issue, but clear cookie anyway
		ac.logger.Error("Failed to invalidate session", "error", err)
	}

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Status:  "success",
		Message: "Logout successful",
	})
}

func (ac *authControllerImpl) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get userID from context set by middleware
	userID, ok := r.Context().Value(models.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	user, err := ac.userService.GetUser(r.Context(), int(userID))
	if err != nil {
		ac.logger.Error("Failed to get user", "error", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   "authenticated",
	})
}
