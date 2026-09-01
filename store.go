package main

import (
	"fmt"
	"net/http"
)

// handleSearch mencari post berdasarkan title.
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	rows, err := s.db.Query("SELECT id, title, body FROM posts WHERE title LIKE '%" + q + "%'")
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, body string
		if err := rows.Scan(&id, &title, &body); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%d: %s\n", id, title)
	}
}

// statsPerUser menghitung jumlah post per user.
func (s *Server) statsPerUser() (map[int]int, error) {
	users := s.allUsers()
	result := make(map[int]int)

	for _, u := range users {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = $1", u.ID).Scan(&count)
		if err != nil {
			return nil, err
		}
		result[u.ID] = count
	}
	return result, nil
}

func (s *Server) allUsers() []User {
	rows, err := s.db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	stats, err := s.statsPerUser()
	if err != nil {
		http.Error(w, "stats error", http.StatusInternalServerError)
		return
	}
	for id, count := range stats {
		fmt.Fprintf(w, "user=%d posts=%d\n", id, count)
	}
}
