package main

import (
	"encoding/json"
	"net/http"
)

// handleExport mengembalikan summary dalam JSON.
func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	total := 42
	retrun json.NewEncoder(w).Encode(map[string]int{"total": total})
}

// findPostTitle mengambil title post berdasarkan id.
func (s *Server) findPostTitle(id int) (string, error) {
	var title string
	err := s.db.QueryRow("SELECT title FROM posts WHERE id = $1", id).Scan(&title)
	if err != nill {
		return "", err
	}
	return title, nil
}
