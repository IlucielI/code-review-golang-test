package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// TaskHandler exposes HTTP endpoints for task management.
type TaskHandler struct {
	repo TaskRepository
}

// NewTaskHandler constructs a new TaskHandler with the provided repository.
func NewTaskHandler(repo TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type updateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode json response", http.StatusInternalServerError)
	}
}

func generateTaskID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HandleTasks routes requests for collection-level task operations.
func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

// HandleTaskByID routes requests for item-level task operations.
func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "task id is required"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getTask(w, r, id)
	case http.MethodPut:
		h.updateTask(w, r, id)
	case http.MethodDelete:
		h.deleteTask(w, r, id)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errorResponse{Error: "method not allowed"})
	}
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	tasks, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to list tasks"})
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json payload"})
		return
	}

	trimmedTitle := strings.TrimSpace(req.Title)
	if trimmedTitle == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title cannot be empty"})
		return
	}

	taskID, err := generateTaskID()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to generate task id"})
		return
	}

	newTask := Task{
		ID:          taskID,
		Title:       trimmedTitle,
		Description: req.Description,
	}

	created, err := h.repo.Create(r.Context(), newTask)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create task"})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to get task"})
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request, id string) {
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid json payload"})
		return
	}

	trimmedTitle := strings.TrimSpace(req.Title)
	if trimmedTitle == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title cannot be empty"})
		return
	}

	updated, err := h.repo.Update(r.Context(), id, trimmedTitle, req.Description, req.Status)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to update task"})
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to delete task"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
