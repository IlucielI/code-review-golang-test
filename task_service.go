package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// TaskService orchestrates business operations on tasks.
type TaskService interface {
	CreateTask(ctx context.Context, title, description string) (Task, error)
	GetTask(ctx context.Context, id string) (Task, error)
	ListTasks(ctx context.Context, limit, offset int) ([]Task, error)
	UpdateTask(ctx context.Context, id string, title, description string, status TaskStatus) (Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type defaultTaskService struct {
	repo TaskRepository
}

// NewTaskService creates a business logic service layer for tasks.
func NewTaskService(repo TaskRepository) TaskService {
	return &defaultTaskService{repo: repo}
}

func (s *defaultTaskService) generateID() (string, error) {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *defaultTaskService) CreateTask(ctx context.Context, title, description string) (Task, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return Task{}, ErrInvalidTitle
	}

	taskID, err := s.generateID()
	if err != nil {
		return Task{}, err
	}

	task := Task{
		ID:          taskID,
		Title:       trimmedTitle,
		Description: strings.TrimSpace(description),
	}

	return s.repo.Create(ctx, task)
}

func (s *defaultTaskService) GetTask(ctx context.Context, id string) (Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *defaultTaskService) ListTasks(ctx context.Context, limit, offset int) ([]Task, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *defaultTaskService) UpdateTask(ctx context.Context, id string, title, description string, status TaskStatus) (Task, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return Task{}, ErrInvalidTitle
	}
	return s.repo.Update(ctx, id, trimmedTitle, strings.TrimSpace(description), status)
}

func (s *defaultTaskService) DeleteTask(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
