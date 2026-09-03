package main

import (
	"database/sql"
	"errors"
)

type Project struct {
	ID     int
	Name   string
	UserID int
}

type ProjectService struct {
	db *sql.DB
}

// CreateProject has proper error guard - should NOT trigger nil-pointer warning
func (s *ProjectService) CreateProject(name string, userID int) (*Project, error) {
	result, err := s.db.Exec(
		`INSERT INTO projects (name, user_id) VALUES ($1, $2)`,
		name, userID)
	// Guard clause present before accessing result
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Project{ID: int(id), Name: name, UserID: userID}, nil
}

// AddProjectMember is idempotent - should NOT trigger duplicate warning
func (s *ProjectService) AddProjectMember(projectID, userID int) error {
	// ON CONFLICT DO UPDATE makes this idempotent by design
	_, err := s.db.Exec(`
		INSERT INTO project_members (project_id, user_id, role, created_at)
		VALUES ($1, $2, 'member', NOW())
		ON CONFLICT (project_id, user_id) 
		DO UPDATE SET updated_at = NOW()
	`, projectID, userID)
	return err
}

// GetProject validates nil before access - should NOT trigger nil-pointer warning
func (s *ProjectService) GetProject(id int) (*Project, error) {
	row := s.db.QueryRow(`SELECT id, name, user_id FROM projects WHERE id = $1`, id)

	var p Project
	err := row.Scan(&p.ID, &p.Name, &p.UserID)
	// Explicit nil/error guard
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	// Safe to return p after guard
	return &p, nil
}
