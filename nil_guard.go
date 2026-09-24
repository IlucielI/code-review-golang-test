package main

import (
	"context"
	"errors"
)

// Test case for false-positive #3: nil-pointer false alarm with guard clause
// This should NOT be flagged because there's explicit guard before dereference

type UserService struct{}

type User struct {
	ID   int
	Name string
}

// CreateUser has proper error guard before dereferencing result
func (s *UserService) CreateUser(ctx context.Context, name string) (*User, error) {
	user, err := s.saveToDatabase(ctx, name)
	// Guard clause: checks error before dereferencing
	if err != nil {
		return nil, err
	}
	// Safe dereference after guard
	return user, nil
}

// GetUserProfile has nil check guard
func (s *UserService) GetUserProfile(ctx context.Context, id int) (*User, error) {
	user, err := s.fetchUser(ctx, id)
	if err != nil {
		return nil, err
	}
	// Another guard for nil pointer
	if user == nil {
		return nil, errors.New("user not found")
	}
	// Safe to access user.Name after guards
	_ = user.Name
	return user, nil
}

// UpdateUserName validates input at handler level
func (s *UserService) UpdateUserName(ctx context.Context, user *User, newName string) error {
	// Request validation happens at handler - if this is called, user is non-nil
	user.Name = newName
	return s.saveToDatabase(ctx, newName)
}

func (s *UserService) saveToDatabase(ctx context.Context, name string) (*User, error) {
	// mock
	return &User{ID: 1, Name: name}, nil
}

func (s *UserService) fetchUser(ctx context.Context, id int) (*User, error) {
	// mock
	return &User{ID: id, Name: "test"}, nil
}
