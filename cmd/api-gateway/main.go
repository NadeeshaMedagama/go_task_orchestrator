package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"go_task_orchestrator/internal/platform/jwtutil"
)

func main() {
	port := getenv("PORT", "8080")
	secret := getenv("JWT_SECRET", "change-me-in-production")
	authURL := mustParseURL(getenv("AUTH_SERVICE_URL", "http://localhost:8081"))
	taskURL := mustParseURL(getenv("TASK_SERVICE_URL", "http://localhost:8082"))

	authProxy := reverseProxy(authURL)
	taskProxy := reverseProxy(taskURL)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","service":"api-gateway"}`))
	})

	mux.Handle("/api/auth/", stripPrefix("/api", authProxy))
	mux.Handle("/api/users", withJWT(secret, stripPrefix("/api", authProxy)))
	mux.Handle("/api/tasks", withJWT(secret, stripPrefix("/api", taskProxy)))
	mux.Handle("/api/tasks/", withJWT(secret, stripPrefix("/api", taskProxy)))

	log.Printf("api-gateway listening on :%s", port)
	if err := http.ListenAndServe(":"+port, logging(mux)); err != nil {
		log.Fatal(err)
	}
}

func reverseProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	originalDirector := proxy.Director
	proxy.Director = func(r *http.Request) {
		originalDirector(r)
		r.Host = target.Host
	}
	return proxy
}

func withJWT(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		claims, err := jwtutil.ParseToken(secret, parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-Name", claims.Username)
		r.Header.Set("X-User-Role", claims.Role)
		next.ServeHTTP(w, r)
	})
}

func stripPrefix(prefix string, next http.Handler) http.Handler {
	return http.StripPrefix(prefix, next)
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatalf("invalid URL %s: %v", raw, err)
	}
	return u
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
