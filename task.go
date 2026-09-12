package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// ErrTaskNotFound is returned when a requested task does not exist.
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidTitle is returned when a task title is empty.
	ErrInvalidTitle = errors.New("task title cannot be empty")
)

// TaskStatus defines the progression state of a task.
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusCompleted TaskStatus = "completed"
)

// Task represents a todo or project work item.
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TaskRepository specifies storage operations for Task entities.
type TaskRepository interface {
	Create(ctx context.Context, task Task) (Task, error)
	GetByID(ctx context.Context, id string) (Task, error)
	List(ctx context.Context, limit, offset int) ([]Task, error)
	Update(ctx context.Context, id string, title, description string, status TaskStatus) (Task, error)
	Delete(ctx context.Context, id string) error
}

// MemoryTaskRepository provides a thread-safe in-memory implementation of TaskRepository.
type MemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

// NewMemoryTaskRepository instantiates a new empty in-memory task repository.
func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{
		tasks: make(map[string]Task),
	}
}

// Create inserts a new task into memory after validating requirements.
func (r *MemoryTaskRepository) Create(ctx context.Context, task Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	if task.Title == "" {
		return Task{}, ErrInvalidTitle
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now
	if task.Status == "" {
		task.Status = StatusPending
	}

	r.tasks[task.ID] = task
	return task, nil
}

// GetByID retrieves a task by its unique identifier.
func (r *MemoryTaskRepository) GetByID(ctx context.Context, id string) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

// List returns a paginated slice of tasks with bounds checking.
func (r *MemoryTaskRepository) List(ctx context.Context, limit, offset int) ([]Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	all := make([]Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		all = append(all, task)
	}

	if offset >= len(all) {
		return []Task{}, nil
	}

	end := offset + limit
	if end > len(all) {
		end = len(all)
	}

	return all[offset:end], nil
}

// Update modifies title, description, and status of an existing task.
func (r *MemoryTaskRepository) Update(ctx context.Context, id string, title, description string, status TaskStatus) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	if title == "" {
		return Task{}, ErrInvalidTitle
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}

	task.Title = title
	task.Description = description
	if status != "" {
		task.Status = status
	}
	task.UpdatedAt = time.Now().UTC()

	r.tasks[id] = task
	return task, nil
}

// Delete removes a task by its identifier.
func (r *MemoryTaskRepository) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}
