package main

import (
	"database/sql"
	"net/http"
)

// Test case for false-positive #5: idempotent upsert operations
// This should NOT be flagged as "duplicate data risk"

type MemberStore struct {
	DB *sql.DB
}

type Member struct {
	ProjectID int
	UserID    int
	Role      string
}

// AddMember is idempotent: ON CONFLICT DO UPDATE means re-accept doesn't duplicate
func (s *MemberStore) AddMember(projectID, userID int, role string) error {
	_, err := s.DB.Exec(`
		INSERT INTO project_members (project_id, user_id, role, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (project_id, user_id) 
		DO UPDATE SET role = EXCLUDED.role, updated_at = NOW()
	`, projectID, userID, role)
	return err
}

// UpsertInvite is also idempotent
func (s *MemberStore) UpsertInvite(email string, projectID int) error {
	_, err := s.DB.Exec(`
		INSERT INTO invites (email, project_id) 
		VALUES ($1, $2)
		ON CONFLICT (email, project_id) DO NOTHING
	`, email, projectID)
	return err
}

// Test case for false-positive #4: public health endpoint
// This should NOT be flagged as "missing auth"

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Public health endpoint - intentionally no auth for monitoring systems
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	// Kubernetes liveness probe endpoint - must be public
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}
