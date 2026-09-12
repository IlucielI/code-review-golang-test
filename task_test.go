package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

func TestTaskService_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryTaskRepository()
	service := NewTaskService(repo)

	// 1. Create task
	created, err := service.CreateTask(ctx, "Clean Code Review Test", "Testing idiomatic clean Go code")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Status != StatusPending {
		t.Errorf("expected pending status, got %s", created.Status)
	}

	// 2. Get task
	fetched, err := service.GetTask(ctx, created.ID)
	if err != nil {
		t.Fatalf("expected task to exist, got %v", err)
	}
	if fetched.Title != created.Title {
		t.Errorf("expected title %s, got %s", created.Title, fetched.Title)
	}

	// 3. Update task
	updated, err := service.UpdateTask(ctx, created.ID, "Clean Code Review Verified", "All tests pass", StatusCompleted)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if updated.Status != StatusCompleted {
		t.Errorf("expected completed status, got %s", updated.Status)
	}

	// 4. List tasks
	list, err := service.ListTasks(ctx, 10, 0)
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 task in list, got %d", len(list))
	}

	// 5. Delete task
	if deleteErr := service.DeleteTask(ctx, created.ID); deleteErr != nil {
		t.Fatalf("expected delete to succeed, got %v", deleteErr)
	}

	// 6. Verify not found after delete
	_, getErr := service.GetTask(ctx, created.ID)
	if getErr != ErrTaskNotFound {
		t.Errorf("expected ErrTaskNotFound, got %v", getErr)
	}
}

func TestTaskService_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	repo := NewMemoryTaskRepository()
	service := NewTaskService(repo)

	var wg sync.WaitGroup
	const numWorkers = 20

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			title := "Task " + strconv.Itoa(workerID)
			item, createErr := service.CreateTask(ctx, title, "Concurrent workload")
			if createErr == nil {
				if _, getErr := service.GetTask(ctx, item.ID); getErr != nil {
					return
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestTaskHandler_HTTPFlow(t *testing.T) {
	repo := NewMemoryTaskRepository()
	service := NewTaskService(repo)
	handler := NewTaskHandler(service)

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
	if decodeErr := json.NewDecoder(rec.Body).Decode(&created); decodeErr != nil {
		t.Fatalf("failed to decode response: %v", decodeErr)
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
