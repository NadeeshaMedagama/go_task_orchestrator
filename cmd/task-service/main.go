package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     string    `json:"dueDate,omitempty"`
	Owner       string    `json:"owner"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"dueDate"`
	Owner       string `json:"owner"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

type pagedResponse struct {
	Content []Task `json:"content"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Total   int    `json:"totalElements"`
}

type server struct {
	mu     sync.RWMutex
	tasks  map[int64]Task
	nextID int64
}

var validStatuses = []string{"TODO", "IN_PROGRESS", "DONE"}
var validPriorities = []string{"LOW", "MEDIUM", "HIGH"}

func main() {
	port := getenv("PORT", "8082")
	s := &server{tasks: make(map[int64]Task), nextID: 1}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("GET /tasks", s.listTasks)
	mux.HandleFunc("POST /tasks", s.createTask)
	mux.HandleFunc("GET /tasks/", s.getTask)
	mux.HandleFunc("PUT /tasks/", s.updateTask)
	mux.HandleFunc("DELETE /tasks/", s.deleteTask)
	mux.HandleFunc("PATCH /tasks/", s.patchTaskStatus)

	log.Printf("task-service listening on :%s", port)
	if err := http.ListenAndServe(":"+port, logging(mux)); err != nil {
		log.Fatal(err)
	}
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "task-service"})
}

func (s *server) listTasks(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-Role")
	username := r.Header.Get("X-User-Name")
	if role == "" || username == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing user identity"})
		return
	}

	statusFilter := strings.ToUpper(r.URL.Query().Get("status"))
	priorityFilter := strings.ToUpper(r.URL.Query().Get("priority"))
	page := intFromQuery(r, "page", 0)
	size := intFromQuery(r, "size", 10)
	if size <= 0 {
		size = 10
	}

	s.mu.RLock()
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if role != "ADMIN" && t.Owner != username {
			continue
		}
		if statusFilter != "" && t.Status != statusFilter {
			continue
		}
		if priorityFilter != "" && t.Priority != priorityFilter {
			continue
		}
		list = append(list, t)
	}
	s.mu.RUnlock()

	slices.SortFunc(list, func(a, b Task) int {
		if a.CreatedAt.Before(b.CreatedAt) {
			return -1
		}
		if a.CreatedAt.After(b.CreatedAt) {
			return 1
		}
		return 0
	})

	start := page * size
	if start > len(list) {
		start = len(list)
	}
	end := start + size
	if end > len(list) {
		end = len(list)
	}

	writeJSON(w, http.StatusOK, pagedResponse{
		Content: list[start:end],
		Page:    page,
		Size:    size,
		Total:   len(list),
	})
}

func (s *server) createTask(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("X-User-Role")
	username := r.Header.Get("X-User-Name")
	if role == "" || username == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing user identity"})
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}
	if req.Status == "" {
		req.Status = "TODO"
	}
	if req.Priority == "" {
		req.Priority = "MEDIUM"
	}

	req.Status = strings.ToUpper(req.Status)
	req.Priority = strings.ToUpper(req.Priority)
	if !slices.Contains(validStatuses, req.Status) || !slices.Contains(validPriorities, req.Priority) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid status or priority"})
		return
	}

	owner := username
	if role == "ADMIN" && req.Owner != "" {
		owner = req.Owner
	}

	now := time.Now()
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	t := Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Owner:       owner,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[id] = t
	s.mu.Unlock()

	writeJSON(w, http.StatusCreated, t)
}

func (s *server) getTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	id, ok := extractTaskID(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}
	t, found := s.lookupTask(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	if !canAccessTask(r, t) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (s *server) updateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		return
	}
	id, ok := extractTaskID(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}
	existing, found := s.lookupTask(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	if !canAccessTask(r, existing) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not allowed"})
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}
	req.Status = strings.ToUpper(req.Status)
	req.Priority = strings.ToUpper(req.Priority)
	if !slices.Contains(validStatuses, req.Status) || !slices.Contains(validPriorities, req.Priority) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid status or priority"})
		return
	}

	existing.Title = req.Title
	existing.Description = req.Description
	existing.Status = req.Status
	existing.Priority = req.Priority
	existing.DueDate = req.DueDate
	existing.UpdatedAt = time.Now()

	s.mu.Lock()
	s.tasks[id] = existing
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, existing)
}

func (s *server) deleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		return
	}
	id, ok := extractTaskID(r.URL.Path)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}
	t, found := s.lookupTask(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	if !canAccessTask(r, t) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not allowed"})
		return
	}

	s.mu.Lock()
	delete(s.tasks, id)
	s.mu.Unlock()
	writeJSON(w, http.StatusNoContent, nil)
}

func (s *server) patchTaskStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch || !strings.HasSuffix(r.URL.Path, "/status") {
		return
	}

	id, ok := extractTaskID(strings.TrimSuffix(r.URL.Path, "/status"))
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task id"})
		return
	}
	t, found := s.lookupTask(id)
	if !found {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	if !canAccessTask(r, t) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "not allowed"})
		return
	}

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	req.Status = strings.ToUpper(req.Status)
	if !slices.Contains(validStatuses, req.Status) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid status"})
		return
	}

	t.Status = req.Status
	t.UpdatedAt = time.Now()
	s.mu.Lock()
	s.tasks[id] = t
	s.mu.Unlock()

	writeJSON(w, http.StatusOK, t)
}

func (s *server) lookupTask(id int64) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func canAccessTask(r *http.Request, t Task) bool {
	role := r.Header.Get("X-User-Role")
	username := r.Header.Get("X-User-Name")
	return role == "ADMIN" || t.Owner == username
}

func extractTaskID(path string) (int64, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 2 || parts[0] != "tasks" {
		return 0, false
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, false
	}
	return id, true
}

func intFromQuery(r *http.Request, key string, fallback int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
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
	if status == http.StatusNoContent {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func getenv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
