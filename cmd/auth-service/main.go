package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go_task_orchestrator/internal/platform/jwtutil"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"tokenType"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type server struct {
	mu      sync.RWMutex
	users   map[string]User
	jwtKey  string
	tokenTT time.Duration
}

func main() {
	port := getenv("PORT", "8081")
	secret := getenv("JWT_SECRET", "change-me-in-production")

	s := &server{
		users: map[string]User{
			"admin": {
				Username: "admin",
				Email:    "admin@go-task-orchestrator.local",
				Password: "admin123",
				Role:     "ADMIN",
			},
		},
		jwtKey:  secret,
		tokenTT: 24 * time.Hour,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /auth/register", s.register)
	mux.HandleFunc("POST /auth/login", s.login)
	mux.HandleFunc("GET /users", s.listUsers)

	log.Printf("auth-service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, logging(mux)); err != nil {
		log.Fatal(err)
	}
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "auth-service"})
}

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "username, email and password are required"})
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[req.Username]; exists {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "username already exists"})
		return
	}

	u := User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "USER",
	}
	s.users[req.Username] = u
	writeJSON(w, http.StatusCreated, map[string]string{
		"username": u.Username,
		"email":    u.Email,
		"role":     u.Role,
	})
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	s.mu.RLock()
	u, ok := s.users[req.Username]
	s.mu.RUnlock()
	if !ok || u.Password != req.Password {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}

	token, err := jwtutil.GenerateToken(s.jwtKey, u.Username, u.Role, s.tokenTT)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		Token:     token,
		TokenType: "Bearer",
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
	})
}

func (s *server) listUsers(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User-Role") != "ADMIN" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "admin access required"})
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]map[string]string, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, map[string]string{
			"username": u.Username,
			"email":    u.Email,
			"role":     u.Role,
		})
	}
	writeJSON(w, http.StatusOK, users)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
