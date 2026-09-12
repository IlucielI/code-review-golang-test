package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// TaskHandler exposes HTTP endpoints for task operations.
type TaskHandler struct {
	service TaskService
}

// NewTaskHandler constructs a new TaskHandler with the provided service.
func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{service: service}
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
	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, writeErr := w.Write(buf); writeErr != nil {
		return
	}
}

func parseQueryInt(param string, defaultVal int) int {
	if param == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return defaultVal
	}
	return val
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
	limit := parseQueryInt(query.Get("limit"), 20)
	offset := parseQueryInt(query.Get("offset"), 0)

	tasks, err := h.service.ListTasks(r.Context(), limit, offset)
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

	created, err := h.service.CreateTask(r.Context(), req.Title, req.Description)
	if err != nil {
		if errors.Is(err, ErrInvalidTitle) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title cannot be empty"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create task"})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.service.GetTask(r.Context(), id)
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

	updated, err := h.service.UpdateTask(r.Context(), id, req.Title, req.Description, req.Status)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, errorResponse{Error: "task not found"})
			return
		}
		if errors.Is(err, ErrInvalidTitle) {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "title cannot be empty"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to update task"})
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteTask(r.Context(), id)
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
