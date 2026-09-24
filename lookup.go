package main

import "database/sql"

// findUserByEmail returns a user by email, nil when not found.
func (s *Server) findUserByEmail(email string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, name FROM users WHERE email = $1", email).Scan(&u.ID, &u.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// getUserName returns the display name for an email.
func (s *Server) getUserName(email string) string {
	user, _ := s.findUserByEmail(email)
	return user.Name
}
