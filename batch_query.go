package main

import (
	"context"
	"database/sql"
)

// Test case for false-positive #9: N+1 batch query detection
// This should NOT be flagged as N+1 because it uses PostgreSQL ANY($1) batch semantics

type UserContact struct {
	ID    string
	Email string
	Phone string
}

type ContactStore struct {
	Pool *sql.DB
}

// ListUserContacts retrieves contacts for multiple user IDs in a single batch query
// This uses ANY($1) which is PostgreSQL's batch WHERE-IN operator
func (s *ContactStore) ListUserContacts(ctx context.Context, userIDs []string) ([]UserContact, error) {
	// Single batch query with ANY($1) - NOT N+1
	rows, err := s.Pool.QueryContext(ctx,
		`SELECT id, email, phone FROM users WHERE id = ANY($1)`, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]UserContact, 0, len(userIDs))
	// This loop scans results, it does NOT query per item
	for rows.Next() {
		var c UserContact
		if err := rows.Scan(&c.ID, &c.Email, &c.Phone); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}

// BatchGetByIDs uses IN clause - also batch query, NOT N+1
func (s *ContactStore) BatchGetByIDs(ctx context.Context, ids []int) ([]UserContact, error) {
	rows, err := s.Pool.QueryContext(ctx,
		`SELECT id, email, phone FROM users WHERE id IN (?)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []UserContact
	for rows.Next() {
		var c UserContact
		rows.Scan(&c.ID, &c.Email, &c.Phone)
		result = append(result, c)
	}
	return result, nil
}
