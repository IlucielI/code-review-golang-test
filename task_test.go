package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestTaskRepository_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryTaskRepository()

	// 1. Create task
	task := Task{
		ID:          "task-1",
		Title:       "Clean Code Review Test",
		Description: "Testing idiomatic clean Go code without defects",
	}
	created, err := repo.Create(ctx, task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Status != StatusPending {
		t.Errorf("expected pending status, got %s", created.Status)
	}

	// 2. Get task
	fetched, err := repo.GetByID(ctx, "task-1")
	if err != nil {
		t.Fatalf("expected task to exist, got %v", err)
	}
	if fetched.Title != task.Title {
		t.Errorf("expected title %s, got %s", task.Title, fetched.Title)
	}

	// 3. Update task
	updated, err := repo.Update(ctx, "task-1", "Clean Code Review Verified", "All tests pass", StatusCompleted)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if updated.Status != StatusCompleted {
		t.Errorf("expected completed status, got %s", updated.Status)
	}

	// 4. List tasks
	list, err := repo.List(ctx, 10, 0)
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 task in list, got %d", len(list))
	}

	// 5. Delete task
	if err := repo.Delete(ctx, "task-1"); err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}

	// 6. Verify not found after delete
	_, err = repo.GetByID(ctx, "task-1")
	if err != ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskRepository_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryTaskRepository()

	var wg sync.WaitGroup
	const numWorkers = 50

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			id := "task-" + string(rune('A'+workerID%26)) + string(rune('0'+workerID/10))
			_, _ = repo.Create(ctx, Task{ID: id, Title: "Concurrent Worker Task"})
			_, _ = repo.GetByID(ctx, id)
			_, _ = repo.List(ctx, 10, 0)
		}(i)
	}
	wg.Wait()
}

func TestTaskHandler_HTTPFlow(t *testing.T) {
	repo := NewMemoryTaskRepository()
	handler := NewTaskHandler(repo)

	// 1. Create task via POST /tasks
	payload := `{"title":"Integration Task","description":"Testing HTTP handlers"}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.HandleTasks(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rec.Code)
	}

	var created Task
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if created.ID == "" {
		t.Errorf("expected generated task ID, got empty string")
	}

	// 2. Get task via GET /tasks/{id}
	req = httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	rec = httptest.NewRecorder()
	handler.HandleTaskByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// 3. Validation error: empty title
	emptyPayload := `{"title":""}`
	req = httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(emptyPayload))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	handler.HandleTasks(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", rec.Code)
	}

	// 4. Not found error
	req = httptest.NewRequest(http.MethodGet, "/tasks/non-existent-id", nil)
	rec = httptest.NewRecorder()
	handler.HandleTaskByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", rec.Code)
	}
}
