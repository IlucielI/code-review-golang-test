package main

import (
	"context"
	"database/sql"
)

type User struct {
	ID       int
	Username string
	Email    string
}

type UserRepository struct {
	db *sql.DB
}

// GetUsersByIDs - Batch query using ANY($1) - should NOT trigger N+1 warning
func (r *UserRepository) GetUsersByIDs(ctx context.Context, userIDs []int) ([]*User, error) {
	// Single batch query with PostgreSQL ANY operator
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, username, email FROM users WHERE id = ANY($1)`, userIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	// Loop over results - this is NOT N+1, it's scanning batch query results
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}

// FindByEmails - Batch query using IN clause - should NOT trigger N+1 warning
func (r *UserRepository) FindByEmails(ctx context.Context, emails []string) ([]*User, error) {
	// Another batch query pattern with IN operator
	query := `SELECT id, username, email FROM users WHERE email IN (?)`
	rows, err := r.db.QueryContext(ctx, query, emails)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Username, &u.Email)
		users = append(users, &u)
	}
	return users, nil
}
